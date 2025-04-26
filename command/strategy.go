package command

import (
    "fmt"
    "maps"
    "os"
    "slices"
    "strings"
    "sync"
    "time"
    "twig/branch"
    "twig/common"
    "twig/log"
    "twig/network"
)

const (
    requestLimit    = 5
    itemsPerRequest = 100

    doneStatusId = 3
)

var rate = time.Tick(time.Second / time.Duration(requestLimit))

type TwigCommand interface {
    Execute() error
}

type createLocalBranchStrategy struct {
    client network.Client
    input  *common.Input
}

func NewCreateLocalBranchCommand(client network.Client, input *common.Input) TwigCommand {
    return &createLocalBranchStrategy{
        client: client,
        input:  input,
    }
}

func (clb *createLocalBranchStrategy) Execute() error {
    issue, ok := clb.input.Arguments[common.Issue]
    if !ok {
        return fmt.Errorf("create command: issue-key must not be null")
    }

    jiraIssue, err := clb.client.GetJiraIssue(issue)
    if err != nil {
        return err
    }

    branchType, err := parseBranchType(*clb.input)
    if err != nil {
        return err
    }

    if branchType == branch.NULL {
        jiraIssueTypes, err := clb.client.GetJiraIssueTypes()
        if err != nil {
            return err
        }

        branchType, err = convertIssueTypeToBranchType(*jiraIssue.Fields.Type, jiraIssueTypes)
        if err != nil {
            return err
        }
    }

    excludePhrases := os.Getenv("branch.exclude")
    if excludePhrases == "" {
        log.Warn().Println("branch.exclude is not set")
    }

    branchName := branch.BuildName(branchType, *jiraIssue, excludePhrases)
    hasBranch := HasBranch(branchName)

    checkoutCommand, err := Checkout(branchName, hasBranch)
    if err != nil {
        return err
    }

    log.Info().Println(checkoutCommand)
    return nil
}

func parseBranchType(input common.Input) (branch.Type, error) {
    branchType, ok := input.Arguments[common.BranchType]
    if !ok {
        log.Debug().Println("get issue type: no user override, take Issue types from Jira")
        return branch.NULL, nil
    }

    log.Debug().Printf("get issue type: user override: %s", branchType)
    return common.ConvertUserInputToBranchType(branchType)
}

func convertIssueTypeToBranchType(jiraIssueType network.IssueType, networkTypes []network.IssueType) (branch.Type, error) {
    localTypes := os.Getenv("branch.mapping")
    if localTypes == "" {
        return branch.NULL, fmt.Errorf("get issue type: branch.mapping is not set")
    }

    mappedIssueTypes, err := common.ConvertIssueTypesToMap(localTypes, networkTypes)
    if err != nil {
        return branch.NULL, fmt.Errorf("get issue type: %w", err)
    }

    value, ok := mappedIssueTypes[jiraIssueType.Id]
    if !ok {
        return branch.NULL, fmt.Errorf("get issue type: mapped issue type does not exist")
    }

    return value, nil
}

type deleteLocalBranchStrategy struct {
    client network.Client
    input  *common.Input
}

func NewDeleteLocalBranchCommand(client network.Client, input *common.Input) TwigCommand {
    return &deleteLocalBranchStrategy{
        client: client,
        input:  input,
    }
}

func (dlb *deleteLocalBranchStrategy) Execute() error {
    fetchCommand, err := ExecuteFetchPrune()
    if err != nil {
        return err
    }
    if fetchCommand != "" {
        log.Info().Println(fetchCommand)
    }

    if err := BranchStatus(); err != nil {
        return err
    }

    devBranch := os.Getenv("branch.default")
    hasBranch := HasBranch(devBranch)

    checkoutCommand, err := Checkout(devBranch, hasBranch)
    if err != nil {
        return err
    }
    log.Info().Println(checkoutCommand)

    localBranches, err := GetLocalBranches()
    if err != nil {
        return err
    }

    issues, err := pairBranchesWithIssues(localBranches)
    if err != nil {
        return err
    }

    statuses, err := pairBranchesWithStatuses(*dlb.input, dlb.client, issues)
    if err != nil {
        return err
    }

    if err := deleteBranchesIfAny(*dlb.input, statuses); err != nil {
        return err
    }

    return nil
}

func deleteBranchesIfAny(input common.Input, statuses map[string]network.IssueStatusCategory) error {
    anyCompleted := false
    remote, hasRemote := input.Arguments[common.Remote]

    for branchName, status := range statuses {
        if status.Id == doneStatusId {
            deleteCommand, err := DeleteLocalBranch(branchName)
            if err != nil {
                log.Error().Print(deleteCommand)
                log.Error().Print(fmt.Errorf("delete local branch: [%s] %w\n", branchName, err))
            } else {
                log.Info().Print(deleteCommand)
            }

            if hasRemote {
                remoteDeleteCommand, err := DeleteRemoteBranch(remote, branchName)
                if err != nil {
                    log.Error().Print(remoteDeleteCommand)
                    log.Error().Print(fmt.Errorf("delete remote branch: [%s] %w\n", branchName, err))
                } else {
                    log.Info().Print(remoteDeleteCommand)
                }
            }

            anyCompleted = true
        }
    }

    if !anyCompleted {
        return fmt.Errorf("delete branch: no associated Jira issues in DONE status")
    }

    return nil
}

func pairBranchesWithStatuses(input common.Input, client network.Client, issues map[string]string) (map[string]network.IssueStatusCategory, error) {
    statuses := make(map[string]network.IssueStatusCategory)
    assignee, hasAssignee := input.Arguments[common.Assignee]

    size := len(issues)

    if size < itemsPerRequest {
        for localBranch, issue := range issues {
            jiraIssue, err := client.GetJiraIssueStatus(issue, hasAssignee)

            if err != nil {
                log.Warn().Printf("pair branch with status: %v", err)
                continue
            }

            if hasAssignee {
                email := jiraIssue.Fields.Assignee.Email

                if err = validateJiraIssue(jiraIssue.Key, email, assignee); err != nil {
                    log.Warn().Printf("pair branch with status: %v", err)
                    continue
                }
            }

            log.Info().Printf("pair branch with status: [%s] : %s", jiraIssue.Fields.Status.Category.Name, localBranch)
            statuses[localBranch] = jiraIssue.Fields.Status.Category
        }
    } else {
        jiraIssues := make([]network.JiraIssue, len(issues))
        values := slices.Collect(maps.Values(issues))
        attemptsNeeded := calculateAttempts(size)

        var wg sync.WaitGroup
        wg.Add(attemptsNeeded)

        for i := 0; i < attemptsNeeded; i++ {
            go func(batch int) {
                defer wg.Done()

                jiraIssues = append(jiraIssues, getJiraIssueStatusBulk(batch, client, values, hasAssignee)...)
            }(i)
        }

        wg.Wait()

        jiraKeyToIssueMap := make(map[string]network.JiraIssue)
        for _, jiraIssue := range jiraIssues {
            jiraKeyToIssueMap[jiraIssue.Key] = jiraIssue
        }

        for localBranch, issue := range issues {
            jiraIssue := jiraKeyToIssueMap[issue]

            if hasAssignee {
                email := jiraIssue.Fields.Assignee.Email

                if err := validateJiraIssue(jiraIssue.Key, email, assignee); err != nil {
                    log.Warn().Printf("pair branch with status: %v", err)
                    continue
                }
            }

            log.Info().Printf("pair branch with status: [%s] : %s", jiraIssue.Fields.Status.Category.Name, localBranch)
            statuses[localBranch] = jiraIssue.Fields.Status.Category
        }
    }

    if len(statuses) == 0 {
        return nil, fmt.Errorf("pair branch with status: no Jira issues in DONE status")
    }

    return statuses, nil
}

func getJiraIssueStatusBulk(batch int, client network.Client, values []string, hasAssignee bool) []network.JiraIssue {
    <-rate

    size := len(values)

    start := batch * itemsPerRequest
    end := start + itemsPerRequest
    if end > size {
        end = size
    }

    jiraIssues, err := client.GetJiraIssueStatusBulk(values[start:end-1], hasAssignee)
    if err != nil {
        log.Warn().Printf("pair branch with status: %v", err)
    }

    return jiraIssues
}

func calculateAttempts(size int) int {
    return (size + itemsPerRequest - 1) / itemsPerRequest
}

func validateJiraIssue(issueKey, email, assignee string) error {
    at := strings.Index(email, "@")
    if at == -1 {
        return fmt.Errorf("validate issue: email %s pulled from Jira issue is either invalid or corrupted", email)
    }

    username := strings.TrimSpace(email[:at])
    if assignee != username {
        return fmt.Errorf("validate issue: %s: assignee provided %s, actual %s", issueKey, assignee, username)
    }

    return nil
}

func pairBranchesWithIssues(rawBranches string) (map[string]string, error) {
    localBranches := strings.Split(rawBranches, "\n")
    issues := make(map[string]string)

    for _, localBranch := range localBranches {
        trimmedBranchName := strings.Join(strings.Fields(localBranch), "")

        issue, err := branch.ExtractIssueNameFromBranch(trimmedBranchName)
        if err != nil || issue == "" {
            continue
        }

        log.Info().Printf("pair branch with issue: [%s] : %s", issue, trimmedBranchName)
        issues[trimmedBranchName] = issue
    }

    if len(issues) == 0 {
        return nil, fmt.Errorf("pair branch with issue: no relation to Jira issues found")
    }

    return issues, nil
}

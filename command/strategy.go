package command

import (
    "brcha/branch"
    "brcha/common"
    "brcha/log"
    "brcha/network"
    "fmt"
    "os"
    "strings"
)

type BrchaCommand interface {
    Execute() error
}

type createLocalBranchStrategy struct {
    client network.Client
    input  *common.Input
}

func NewCreateLocalBranchCommand(client network.Client, input *common.Input) BrchaCommand {
    return &createLocalBranchStrategy{
        client: client,
        input:  input,
    }
}

func (clb *createLocalBranchStrategy) Execute() error {
    jiraIssue, err := clb.client.GetJiraIssue(clb.input.Issue)
    if err != nil {
        return err
    }

    branchType, err := parseBranchType(clb.input)
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

    branchName := branch.BuildName(branchType, *jiraIssue)
    hasBranch := HasBranch(branchName)

    checkoutCommand, err := Checkout(branchName, hasBranch)
    if err != nil {
        return err
    }

    log.Info().Println(checkoutCommand)
    return nil
}

func parseBranchType(input *common.Input) (branch.Type, error) {
    if len(input.Argument) > 0 {
        log.Debug().Printf("get issue type: user override: %s", input.Argument)
        return common.ConvertUserInputToBranchType(input.Argument)
    }
    log.Debug().Println("get issue type: no user override, take Issue Types from Jira")

    return branch.NULL, nil
}

func convertIssueTypeToBranchType(jiraIssueType network.IssueType, networkTypes []network.IssueType) (branch.Type, error) {
    localTypes := os.Getenv("BRCHA_MAPPING")

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
}

func NewDeleteLocalBranchCommand(client network.Client) BrchaCommand {
    return &deleteLocalBranchStrategy{
        client: client,
    }
}

func (dlb *deleteLocalBranchStrategy) Execute() error {
    fetchCommand, err := ExecuteFetchPrune()
    if err != nil {
        return err
    }
    log.Info().Println(fetchCommand)

    if err := BranchStatus(); err != nil {
        return err
    }

    devBranch := os.Getenv("BRCHA_DEV_BRANCH_NAME")
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

    statuses, err := pairBranchesWithStatuses(dlb.client, issues)
    if err != nil {
        return err
    }

    deleteCommand, err := deleteBranchesIfAny(statuses)
    if err != nil {
        return err
    }
    log.Info().Printf("delete branches: %s", deleteCommand)

    return nil
}

func deleteBranchesIfAny(statuses map[string]network.IssueStatusCategory) (string, error) {
    log.Info().Printf("delete branches: %+v", statuses)

    var logs string
    for branchName, status := range statuses {
        if status.Id == 3 {
            deleteCommand, err := DeleteLocalBranch(branchName)
            if err != nil {
                return "", err
            }

            logs += "\n" + deleteCommand
        }
    }

    if logs == "" {
        return "", fmt.Errorf("delete branches: no Jira issues in DONE status")
    }

    return logs, nil
}

func pairBranchesWithStatuses(client network.Client, issues map[string]string) (map[string]network.IssueStatusCategory, error) {
    log.Info().Println("pair branch with status")

    statuses := make(map[string]network.IssueStatusCategory)

    for issue, localBranch := range issues {
        jiraIssue, err := client.GetJiraIssueStatus(issue)

        if err != nil {
            continue
        }

        log.Debug().Printf("pair branch with status: %+v %s", jiraIssue.Fields.Status.Category, localBranch)
        statuses[localBranch] = jiraIssue.Fields.Status.Category
    }

    if len(statuses) == 0 {
        return nil, fmt.Errorf("pair branch with status: no Jira issues in DONE status")
    }

    return statuses, nil
}

func pairBranchesWithIssues(rawBranches string) (map[string]string, error) {
    log.Info().Println("pair branch with issue")

    localBranches := strings.Split(rawBranches, "\n")
    issues := make(map[string]string)

    for _, localBranch := range localBranches {
        issue, err := branch.ExtractIssueNameFromBranch(localBranch)
        if err != nil || issue == "" {
            continue
        }

        log.Debug().Printf("pair branch with issue: [%s] %s", issue, localBranch)
        issues[issue] = localBranch
    }

    if len(issues) == 0 {
        return nil, fmt.Errorf("pair branch with issue: no relation to Jira issues found")
    }

    return issues, nil
}

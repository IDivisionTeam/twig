package main

import (
    "brcha/branch"
    "brcha/command"
    "brcha/common"
    "brcha/log"
    "brcha/network"
    "flag"
    "fmt"
    "net/http"
    "os"
)

const (
    emptyCommandArguments string = "Use \"brcha -h\" or \"brcha -help\" for more information."
    helpCommandOutput     string = `
    Usage:
        brcha [arguments]

    The arguments are:
        -i <issue-key>
        -t <branch-type>

    Available branch types:
        build, b: Changes that affect the build system or external dependencies (example scopes: gradle, npm)
        chore, ch: Routine tasks that don't affect the functionality or user-facing aspects of a project
        ci: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
        docs, d: Documentation only changes
        feat, ft: A new feature
        fix, fx: A bug fix
        perf, p: A code change that improves performance
        refactor, rf: A code change that neither fixes a bug nor adds a feature
        revert, rv: A code that restors to a previous or default condition
        style, s: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
        test, t: Adding missing tests or correcting existing tests

    Example:
        ~% brcha -i XX-111 -t fx
        ~% git checkout -b fix/XX-111_jira-issue-name`
)

func main() {
    if err := command.ReadEnvVariables(); err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    input := readUserInput()

    httpClient := &http.Client{}
    client := network.NewClient(httpClient)

    jiraIssue, err := client.GetJiraIssue(input.Issue)
    if err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    branchType, err := parseBranchType(input)
    if err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    if branchType == branch.NULL {
        jiraIssueTypes, err := client.GetJiraIssueTypes()
        if err != nil {
            log.Error().Println(err)
            os.Exit(1)
        }

        branchType, err = convertIssueTypeToBranchType(jiraIssue.Fields.Type, jiraIssueTypes)
        if err != nil {
            log.Error().Println(err)
            os.Exit(1)
        }
    }

    branchName := branch.BuildName(branchType, *jiraIssue)
    hasBranch := command.HasBranch(branchName)

    checkoutCommand, err := command.Checkout(branchName, hasBranch)
    if err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    log.Info().Println(checkoutCommand)
}

func readUserInput() *common.Input {
    var input = &common.Input{
        Issue:    "",
        Argument: "",
    }

    flag.StringVar(&input.Issue, "i", "", "issue key")
    flag.StringVar(&input.Argument, "t", "", "(optional) overrides the type of branch")
    help := flag.Bool("help", false, "displays all available commands")

    flag.Parse()

    if *help == true {
        log.Info().Println(helpCommandOutput)
        os.Exit(0)
    }

    if len(os.Args) == 1 {
        log.Info().Println(emptyCommandArguments)
        os.Exit(0)
    }

    log.Debug().Printf("user input: -i=%s -t=%s", input.Issue, input.Argument)
    return input
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

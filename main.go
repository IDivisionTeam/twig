package main

import (
    "brcha/branch"
    "brcha/command"
    "brcha/common"
    "brcha/network"
    "brcha/recorder"
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
        recorder.Println(recorder.ERROR, err)
        os.Exit(1)
    }

    input := readUserInput()

    httpClient := &http.Client{}
    client := network.NewClient(httpClient)

    jiraIssue, err := client.GetJiraIssue(input.Issue)
    if err != nil {
        recorder.Println(recorder.ERROR, err)
        os.Exit(1)
    }

    jiraIssueTypes, err := client.GetJiraIssueTypes()
    if err != nil {
        recorder.Println(recorder.ERROR, err)
        os.Exit(1)
    }

    branchType, err := getIssueType(input, jiraIssue.Fields.Type, jiraIssueTypes)
    if err != nil {
        recorder.Println(recorder.ERROR, err)
        os.Exit(1)
    }

    branchName := branch.BuildName(branchType, *jiraIssue)

    hasBranch := command.HasBranch(branchName)
    checkoutCommand, err := command.Checkout(branchName, hasBranch)
    if err != nil {
        recorder.Println(recorder.ERROR, err)
        os.Exit(1)
    }

    recorder.Println(recorder.INFO, checkoutCommand)
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
        recorder.Println(recorder.INFO, helpCommandOutput)
        os.Exit(0)
    }

    if len(os.Args) == 1 {
        recorder.Println(recorder.INFO, emptyCommandArguments)
        os.Exit(0)
    }

    return input
}

func getIssueType(input *common.Input, jiraIssueType network.IssueType, types []network.IssueType) (branch.Type, error) {
    if len(input.Argument) > 0 {
        return common.ConvertUserInputToBranchType(input.Argument)
    }

    mappedIssueTypes, err := common.ConvertIssueTypesToMap(types)
    if err != nil {
        return branch.NULL, fmt.Errorf("getIssueType: %w", err)
    }

    value, ok := mappedIssueTypes[jiraIssueType.Id]
    if !ok {
        return branch.NULL, fmt.Errorf("getIssueType: mapped issue type does not exist")
    }

    return value, nil
}

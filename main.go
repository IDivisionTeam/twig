package main

import (
    "brcha/branch"
    "brcha/command"
    "brcha/common"
    "brcha/network"
    "flag"
    "fmt"
    "log"
    "os"
)

const helpCommandOutput string = `
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

func main() {
    err := command.ReadEnvVariables()
    if err != nil {
        log.Fatal(err)
    }

    input, err := readUserInput()
    if err != nil {
        log.Fatal(err)
    }

    jiraIssue, err := network.GetJiraIssue(input.Issue)
    if err != nil {
        log.Fatal(err)
    }

    jiraIssueTypes, err := network.GetJiraIssueTypes()
    if err != nil {
        log.Fatal(err)
    }

    branchType, err := getIssueType(input, jiraIssue.Fields.Type, jiraIssueTypes)
    if err != nil {
        log.Fatal(err)
    }

    branchName := branch.BuildName(branchType, jiraIssue)

    executableCommand, err := command.Checkout(branchName)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(executableCommand)
}

func readUserInput() (*common.Input, error) {
    var input = &common.Input{
        Issue:    "",
        Argument: "",
    }

    help := flag.Bool("help", false, "displays all available commands")
    flag.StringVar(&input.Issue, "i", "", "issue key")
    flag.StringVar(&input.Argument, "t", "", "(optional) overrides the type of branch")

    flag.Parse()

    if *help == true {
        fmt.Println(helpCommandOutput)
        os.Exit(0)
    }

    if len(os.Args) == 1 {
        fmt.Println("Use \"brcha -h\" or \"brcha -help\" for more information.")
        os.Exit(0)
    }

    return input, nil
}

func getIssueType(input *common.Input, jiraIssueType network.IssueType, types []network.IssueType) (string, error) {
    if len(input.Argument) > 0 {
        return common.ConvertUserInputToBranchType(input.Argument)
    }

    mappedIssueTypes, err := common.ConvertIssueTypesToMap(types)
    if err != nil {
        return "", err
    }

    value, ok := mappedIssueTypes[jiraIssueType.Id]
    if !ok {
        return "", fmt.Errorf("mapped issue type does not exist")
    }

    return value, nil
}

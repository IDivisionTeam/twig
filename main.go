package main

import (
    "brcha/branch"
    "brcha/command"
    "brcha/common"
    "brcha/network"
    "fmt"
    "log"
    "os"
)

const helpCommandOutput string = `
Usage: 
    brcha <jira-key> [arguments]

The arguments are: 
    -t  <branch-type>

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
    ~% brcha XX-111 -t fx
    ~% git checkout -b fix/XX-111_jira-issue-name`

func main() {
    input, err := readUserInput()
    if err != nil {
        log.Panic(err)
    }

    jiraIssue, err := network.GetJiraIssue(input.ComandOrIssue)
    if err != nil {
        log.Panic(err)
    }

    jiraIssueTypes, err := network.GetJiraIssueTypes()
    if err != nil {
        log.Panic(err)
    }

    branchType, err := getIssueType(input, jiraIssue.Fields.Type, jiraIssueTypes)
    if err != nil {
        log.Panic(err)
    }

    branchName := branch.BuildName(branchType, jiraIssue)

    executableCommand, err := command.Checkout(branchName)
    if err != nil {
        log.Panic(err)
    }

    fmt.Println(executableCommand)
}

func readUserInput() (*common.Input, error) {
    var input = &common.Input{
        ComandOrIssue: branch.NewBranchType().Chore,
        Argument:      "",
    }

    if len(os.Args) > 3 {
        return input, fmt.Errorf("too many arguments")
    }

    if len(os.Args) == 1 {
        log.Println("Use \"brcha help\" for more information.")
        os.Exit(0)
    }

    if len(os.Args) == 2 {
        input.ComandOrIssue = os.Args[1]

        if input.ComandOrIssue == "help" {
            log.Println(helpCommandOutput)
            os.Exit(0)
        }
    }

    if len(os.Args) == 3 {
        input.ComandOrIssue = os.Args[1]
        arg := os.Args[2]

        if arg != "-t" {
            return input, fmt.Errorf("unsupported argument")
        }

        input.Argument = os.Args[3]
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

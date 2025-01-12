package main

import (
    "brcha/command"
    "brcha/common"
    "brcha/log"
    "brcha/network"
    "flag"
    "net/http"
    "os"
)

const (
    emptyCommandArguments string = `Use "brcha -h" or "brcha -help" for more information.`
    helpCommandOutput     string = `
    Usage:
        brcha [arguments]

    The arguments are:
        -i <issue-key>
        -t <branch-type>
        -clean
        -r <remote>
        -assignee <username>

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

    Examples:
        ~% brcha -i XX-111
        ~% branch created: task/XX-111_jira-issue-name
        ~%
        ~% brcha -i XX-111 -t fx
        ~% branch created: fix/XX-111_jira-issue-name
        ~%
        ~% brcha -clean
        ~% branch deleted: fix/XX-111_jira-issue-name
        ~%
        ~% brcha -clean -r origin
        ~% branch deleted: fix/XX-111_jira-issue-name
        ~% branch deleted: origin/fix/XX-111_jira-issue-name
        ~%
        ~% brcha -clean -r origin -assignee example.user
        ~% branch deleted: fix/XX-111_jira-issue-name
        ~% branch deleted: origin/fix/XX-111_jira-issue-name`
)

func main() {
    input := readUserInput()

    if input.HasFlag(common.HelpFlag) {
        log.Info().Println(helpCommandOutput)
        os.Exit(0)
    }

    if err := command.ReadEnvVariables(); err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    httpClient := &http.Client{}
    client := network.NewClient(httpClient)

    var cmd command.BrchaCommand
    if input.HasFlag(common.CleanFlag) {
        cmd = command.NewDeleteLocalBranchCommand(client, input)
    } else {
        cmd = command.NewCreateLocalBranchCommand(client, input)
    }

    if err := cmd.Execute(); err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }
}

func readUserInput() *common.Input {
    var input = &common.Input{
        Flags:     common.EmptyFlag,
        Arguments: make(map[common.InputType]string),
    }

    help := flag.Bool("help", false, "displays all available commands")
    issue := flag.String("i", "", "issue key")
    branchType := flag.String("t", "", "(optional) overrides the type of branch")
    clean := flag.Bool("clean", false, "deletes all local branches with Jira status Done")
    remote := flag.String("r", "", "(optional) provides remote to delete branch in origin")
    assignee := flag.String("assignee", "", "(optional) provides assignee to delete remote branch")

    flag.Parse()

    if help != nil && *help {
        input.AddFlag(common.HelpFlag)
    }

    if issue != nil && *issue != "" {
        input.Arguments[common.Issue] = *issue

        if branchType != nil && *branchType != "" {
            input.Arguments[common.BranchType] = *branchType
        }
    }

    if clean != nil && *clean {
        input.AddFlag(common.CleanFlag)

        if remote != nil && *remote != "" {
            input.Arguments[common.Remote] = *remote
        }

        if assignee != nil && *assignee != "" {
            input.Arguments[common.Assignee] = *assignee
        }
    }

    if (len(os.Args) == 1) || input.Flags == common.EmptyFlag && (len(input.Arguments) == 0) {
        log.Error().Println(emptyCommandArguments)
        os.Exit(0)
    }

    log.Debug().Printf("user input: flags=%d args=%+v", input.Flags, input.Arguments)
    return input
}

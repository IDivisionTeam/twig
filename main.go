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
    emptyCommandArguments string = `
    Use \"brcha -h\" or \"brcha -help\" for more information.`
    helpCommandOutput string = `
    Usage:
        brcha [arguments]

    The arguments are:
        -i <issue-key>
        -t <branch-type>
        -clean
        -o <remote-origin>

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
        ~% brcha -clean -o origin
        ~% branch deleted: fix/XX-111_jira-issue-name
        ~% branch deleted: origin/fix/XX-111_jira-issue-name`
)

func main() {
    if err := command.ReadEnvVariables(); err != nil {
        log.Error().Println(err)
        os.Exit(1)
    }

    input := readUserInput()

    httpClient := &http.Client{}
    client := network.NewClient(httpClient)

    var cmd command.BrchaCommand
    if input.Issue == "" {
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
        Issue:    "",
        Argument: "",
    }

    flag.StringVar(&input.Issue, "i", "", "issue key")
    flag.StringVar(&input.Argument, "t", "", "(optional) overrides the type of branch")
    flag.StringVar(&input.Argument, "o", "", "(optional) provides origin to delete remote branch")
    help := flag.Bool("help", false, "displays all available commands")
    clean := flag.Bool("clean", false, "deletes all local branches with Jira status Done")

    flag.Parse()

    if *help == true {
        log.Info().Println(helpCommandOutput)
        os.Exit(0)
    }

    if *clean == true {
        log.Debug().Printf("user input: -o=%s", input.Argument)
        return input
    }

    if (len(os.Args) == 1) || (input.Issue == "") {
        log.Info().Println(emptyCommandArguments)
        os.Exit(0)
    }

    log.Debug().Printf("user input: -i=%s -t=%s", input.Issue, input.Argument)
    return input
}

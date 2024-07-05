# brcha | `v1.0.0`

## Overview

A tool for seamless branch creation by passing the Issue Key into the CLI. It uses the Git and Jira APIs under the hood.

## Features

- Seamless Branch Creation: Easily create new Git branches by simply passing the Issue Key into the CLI.
- Integration with Git and Jira: Leverages the Git and Jira APIs to streamline your workflow.
- Automated Branch Naming: Automatically generates branch names based on Jira Ticket Summary, ensuring consistent and
  meaningful branch names.

## Installation

1. Configure your Jira API settings in the `network.go` file.

```network.go
jiraCredentials {
    host:  "example.atlassian.net",
    email: "email@example.com",
    token: "api_token",
}
```

2. Define mappings for your Jira issue types in the configuration the `issue.go` file.

```issue.go
const (
    Build    = "xxxxx"
    Chore    = "xxxxx"
    Ci       = "xxxxx"
    Docs     = "xxxxx"
    Feat     = "xxxxx"
    Fix      = "xxxxx"
    Perf     = "xxxxx"
    Refactor = "xxxxx"
    Revert   = "xxxxx"
    Style    = "xxxxx"
    Test     = "xxxxx"
)
```

You can `curl` available `issuetype`s from Jira.

```terminal
curl --request GET \
  --url 'https://{host}/rest/api/2/issuetype' \
  --user 'email@example.com:<api_token>' \
  --header 'Accept: application/json'
```

3. Specify issue type IDs to ignore, if necessary.

```issue.go
builder.Set("xxxxx", true) // Subtask
builder.Set("xxxxx", true) // Epic
// etc.
```

4. Compile the tool into an executable file.

```terminal
go build
```

5. Move the executable in `/usr/local/bin` for easy global access.

```terminal
mv brcha /usr/local/bin
```

## Usage

```terminal
brcha <jira-key> [arguments]
```

## Commands

`help` - Displays help information for all available commands and options in the CLI tool, providing usage instructions
and examples. Use this command to understand how to use the tool effectively.

```terminal
brcha help
```

`-t <branch-type>` - (optional) Overrides the type of branch to create, allowing the branch name ignore mapped Jira
issue types. Branches are named according to the [standard](https://www.conventionalcommits.org/en/v1.0.0/).

``` terminal
brcha XXX-00 -t ci
```

Available branch types

- `build`, `b` - Changes that affect the build system or external dependencies (example scopes: gradle, npm)
- `chore`, `ch` - Routine tasks that don't affect the functionality or user-facing aspects of a project
- `ci` - Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- `docs`, `d` - Documentation only changes
- `feat`, `ft` - A new feature
- `fix`, `fx` - A bug fix
- `perf`, `p` - A code change that improves performance
- `refactor`, `rf` - A code change that neither fixes a bug nor adds a feature
- `revert`, `rv` - A code that restors to a previous or default condition
- `style`, `s` - Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `test`, `t` - Adding missing tests or correcting existing tests

## Configuration

In case you want to experiment with custom branch formatting or extend existing methods go to `branch.go` file.

```branch.go
func BuildName(branchType string, jiraIssue JiraIssue) string {
    summary := replacePhrases(jiraIssue.Fields.Summary)
    summary = strings.ToLower(summary)
    summary = strings.TrimSpace(summary)
    summary = stripRegex(summary)
    summary = strings.TrimSuffix(summary, "-")

    ...

    // returns "branchType/jiraIssue.Key_summary"
}
```

## Examples

```terminal
~% brcha XX-111 -t fx
~% git checkout -b fix/XX-111_jira-issue-name
```

```terminal
~% brcha XX-111
~% git checkout -b chore/XX-111_jira-issue-name
```
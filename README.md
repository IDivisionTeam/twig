# brcha  ![brcha](https://img.shields.io/badge/brcha-v1.0.1-green.svg)

## Overview

A tool for seamless branch creation by passing the Issue Key into the CLI. It uses the Git and Jira APIs under the hood.

## Features

- Seamless Branch Creation: Easily create new Git branches by simply passing the Issue Key into the CLI.
- Integration with Git and Jira: Leverages the Git and Jira APIs to streamline your workflow.
- Automated Branch Naming: Automatically generates branch names based on Jira Ticket Summary, ensuring consistent and
  meaningful branch names.

## Installation

1. Configure your Jira API settings in the `.env` file.

```.env
BRCHA_HOST=example.atlassian.net
BRCHA_EMAIL=email@example.com
BRCHA_TOKEN=api_token
```
>*NOTE: for* `Bearer` *auth leave* `BRCHA_EMAIL` *field empty!*

2. Define mappings for your Jira issue types in the configuration the `.env` file. Use zero if you want to ignore a specific type.

```.env
BRCHA_MAPPING=build:0;chore:0;ci:0;docs:0;feat:0;fix:0;pref:0;refactor:0;revert:0;style:0;test:0
```

If you have multiple IDs of the same type, separate them with a comma (`,`).
```.env
BRCHA_MAPPING=build:10001,10002,10003;...
```

You can `curl` available `issuetype`s from Jira.

```terminal
curl \
    -D- \
    -X GET \
    -u "email@example.com:token" \
    -H "Content-Type: application/json" \
    https://{host}/rest/api/2/issuetype
```
or
```terminal
curl \
    -D- \
    -X GET \
    -H "Authorization: Bearer {token}" \
    -H "Content-Type: application/json" \
    https://{host}/rest/api/2/issuetype
```

3. Specify the branch that will serve as the base when checking out before deleting local branches. This ensures consistency and avoids issues during cleanup operations.
```.env
BRCHA_DEV_BRANCH_NAME=develop
```

4. Copy `.env` file into `~/.config/brcha/` folder.

```terminal
mkdir -p ~/.config/brcha/ && \
cp .env ~/.config/brcha/
```

5. Compile the tool into an executable file.

```terminal
go build
```

6. Move the executable  into `/usr/local/bin` for easy global access.

```terminal
mv brcha /usr/local/bin
```

## Usage

```terminal
brcha [arguments]
```

## Commands

`help` - Displays help information for all available commands and options in the CLI tool, providing usage instructions
and examples. Use this command to understand how to use the tool effectively.

```terminal
brcha -help
```

`-i <issue-key>` - The branch prefix after branch type. Uses Jira Issue Key.

``` terminal
brcha -i XXX-00
```

`-t <branch-type>` - (optional) Overrides the type of branch to create, allowing the branch name ignore mapped Jira
issue types. Branches are named according to the [standard](https://www.conventionalcommits.org/en/v1.0.0/).

``` terminal
brcha -i XXX-00 -t ci
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

`clean` - Deletes all local branches which have Jira tickets in 'Done' state.

```terminal
brcha -clean
```

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
~% brcha -i XX-111
~% branch created: task/XX-111_jira-issue-name
```

```terminal
~% brcha -i XX-111 -t fx
~% branch created: fix/XX-111_jira-issue-name
```

```terminal
~% brcha -clean
~% branch deleted: fix/XX-111_jira-issue-name`
```
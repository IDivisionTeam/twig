# twig ![tests](https://github.com/yaroslav-android/twig/actions/workflows/go.yml/badge.svg)


## Overview

Streamline your workflow with a CLI tool that integrates Git and Jira. Quickly create, name, and delete branches using the Issue Key, ensuring consistency and efficiency in branch management.

## Features

- Seamless Branch Creation: Easily create new Git branches by simply passing the Issue Key into the CLI.
- Integration with Git and Jira: Leverages the Git and Jira APIs to streamline your workflow.
- Automated Branch Naming: Automatically generates branch names based on Jira Ticket Summary, ensuring consistent and
  meaningful branch names.
- Branch Deletion: Conveniently delete branches directly from the CLI for efficient repository management.

## Installation

1. Configure your Jira API settings in the `.env` file.

```.env
project.host=example.atlassian.net
project.email=example.user@example.com
project.token=api_token
```
>*NOTE: for* `Bearer` *auth leave* `project.email` *field empty!*

2. Define mappings for your Jira issue types in the configuration the `.env` file. Use zero if you want to ignore a specific type.

```.env
branch.mapping=build:0;chore:0;ci:0;docs:0;feat:0;fix:0;pref:0;refactor:0;revert:0;style:0;test:0
```

If you have multiple IDs of the same type, separate them with a comma (`,`).
```.env
branch.mapping=build:10001,10002,10003;...
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
branch.default=develop
```

4. Specify any exclusion phrases to be removed from the branch name, if applicable.
```.env
branch.exclude=front,mobile,android,ios,be,web,spike,eval
```

5. Copy `.env` file into `~/.config/twig/` folder.

```terminal
mkdir -p ~/.config/twig/ && \
cp .env ~/.config/twig/
```

6. Compile the tool into an executable file or [download compiled executable](https://github.com/yaroslav-android/twig/releases).

```terminal
go build
```

7. Move the executable  into `/usr/local/bin` for easy global access.

```terminal
mv twig /usr/local/bin
```

## Usage

```terminal
twig [arguments]
```

## Commands
> _Note: Remote branches can only be deleted if a corresponding local branch exists._

`help` - Displays help information for all available commands and options in the CLI tool, providing usage instructions
and examples. Use this command to understand how to use the tool effectively.

```terminal
btwig -help
```

`i <issue-key>` - The branch prefix after branch type. Uses Jira Issue Key.

``` terminal
twig -i XXX-00
```

`t <branch-type>` - (optional) Overrides the type of branch to create, allowing the branch name ignore mapped Jira
issue types. Branches are named according to the [standard](https://www.conventionalcommits.org/en/v1.0.0/).

``` terminal
twig -i XXX-00 -t ci
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
twig -clean
```

`r <remote>` - (optional) Allows the deletion of remote branches alongside their corresponding local branches.

```terminal
twig -clean -r origin
```

`assignee <username>` - (optional) Specifies the username (from the email) to verify that the Jira issue is assigned to you before permitting remote branch deletion. Use your Jira email, which might match `project.email`, e.g., `example.user@example.com`.

```terminal
twig -clean -r origin -assignee example.user
```

## Configuration

In case you want to experiment with custom branch formatting or extend existing methods go to `branch.go` file.

```branch.go
func BuildName(bt Type, jiraIssue network.JiraIssue, excludePhrases string) string {
    ...

    return "branchType/jiraIssue.Key_summary"
}
```

## Examples

```terminal
~% twig -i XX-111
~% branch created: task/XX-111_jira-issue-name
```

```terminal
~% twig -i XX-111 -t fx
~% branch created: fix/XX-111_jira-issue-name
```

```terminal
~% twig -clean
~% branch deleted: fix/XX-111_jira-issue-name
```

```terminal
~% twig -clean -r origin
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

```terminal
~% twig -clean -r origin -assignee example.user
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

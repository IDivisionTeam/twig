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

1. Configure your Jira API settings in the `twig.config` file.

```twig.config
project.host=example.atlassian.net
project.email=example.user@example.com
project.token=api_token
```

> *NOTE: for* `Bearer` *auth leave* `project.email` *field empty!*

2. Define mappings for your Jira issue types in the configuration the `twig.config` file. Use zero if you want to ignore a specific type.

```twig.config
branch.mapping=build:0;chore:0;ci:0;docs:0;feat:0;fix:0;pref:0;refactor:0;revert:0;style:0;test:0
```

If you have multiple IDs of the same type, separate them with a comma (`,`).

```twig.config
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

3. Specify the `branch.default` which will serve as the base when checking out before deleting local branches. Specify the `branch.origin` to be able to delete branches alongside their corresponding local branches. This ensures consistency and avoids issues during cleanup operations.

```twig.config
branch.default=develop
branch.origin=origin
```

4. Specify any exclusion phrases to be removed from the branch name, if applicable.

```twig.config
branch.exclude=front,mobile,android,ios,be,web,spike,eval
```

5. Copy `twig.config` file into `~/.config/twig/` folder.

```terminal
mkdir -p ~/.config/twig/ && \
cp twig.config ~/.config/twig/
```

6. Compile the tool into an executable file or [download compiled executable](https://github.com/yaroslav-android/twig/releases).

```terminal
go build
```

7. Move the executable into `/usr/local/bin` for easy global access.

```terminal
mv twig /usr/local/bin
```

## Usage

```terminal
twig [-h | --help] [-v | --version] [--config]
```

## Commands

```> twig-help
twig help <command>
```

Displays help information for all available commands and options in the CLI tool providing usage instructions and examples. Use this command to understand how to use the tool effectively.

**Examples**

```terminal
twig help create
```

<br/>

```> twig-create
twig create <issue> [-t | --type]
```

Creates the branch using Jira Issue Key as prefix after branch type.

**Options**

`-t` <br/>
`--type` - (optional) Overrides the type of branch to create, allowing the branch name ignore mapped Jira issue types. Branches are named according to the [standard](https://www.conventionalcommits.org/en/v1.0.0/).

**Available branch types**

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

**Examples**

```terminal
twig create XX-111
```

<br/>

```> twig-clean
twig clean local
twig clean all [-a | --assignee]
```

Deletes branches which have Jira tickets in 'Done' state.<br/>
Note: Remote branches can only be deleted if a corresponding local branch exists.

`-a` <br/>
`--assignee` - (optional) Specifies the username (from the email) to verify that the Jira issue is assigned to you before permitting remote branch deletion. Use your Jira email, which might match `project.email`, e.g., `example.user@example.com`.

**Examples**
```terminal
twig clean local
```
```terminal
twig clean all -a example.user
```

<br/>

## Configuration

In case you want to experiment with custom branch formatting or extend existing methods go to `branch.go` file.

```branch.go
func BuildName(bt Type, jiraIssue network.JiraIssue, excludePhrases string) string {
    ...

    return "branchType/jiraIssue.Key_summary"
}
```

## More Examples

```terminal
~% twig create XX-111
~% branch created: task/XX-111_jira-issue-name
```

```terminal
~% twig create XX-111 -t fx
~% branch created: fix/XX-111_jira-issue-name
```

```terminal
~% twig clean local
~% branch deleted: fix/XX-111_jira-issue-name
```

```terminal
~% twig clean all
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

```terminal
~% twig clean all --assignee example.user
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

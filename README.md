# twig ![tests](https://github.com/yaroslav-android/twig/actions/workflows/go.yml/badge.svg)

<br/>

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [twig-help](#twig-help)
    - [twig-init](#twig-init)
    - [twig-create](#twig-create)
    - [twig-clean](#twig-clean)
- [Configuration](#configuration)
- [More Examples](#more-examples)

<br/>

## Overview

Streamline your workflow with a CLI tool that integrates Git and Jira. Quickly create, name, and delete branches using the Issue Key, ensuring consistency and efficiency in branch management.

## Features

- Seamless Branch Creation: Easily create new Git branches by simply passing the Issue Key into the CLI.
- Integration with Git and Jira: Leverages the Git and Jira APIs to streamline your workflow.
- Automated Branch Naming: Automatically generates branch names based on Jira Ticket Summary, ensuring consistent and
  meaningful branch names.
- Branch Deletion: Conveniently delete branches directly from the CLI for efficient repository management.

## Installation
> Preferably, start from the step 6 and use the [twig-init](#twig-init) command to set up your configuration.

1. Configure your Jira API and VCS settings in the `twig/config/twig.toml` file.

```
[project]
host  = "example.atlassian.net"
auth = "basic"
email = "example.user@example.com"
token = "api_token"
```

> *NOTE: for* `Bearer` *auth use* `bearer` *key for auth property!*

2. Define mappings for your Jira issue types in the configuration the `twig.toml` file. Use **zero** if you want to ignore a specific type.

```
[mapping]
build = ["0"]
chore = ["0"]
ci = ["0"]
docs = ["0"]
feat = ["0"]
fix = ["0"]
pref = ["0"]
refactor = ["0"]
revert = ["0"]
style = ["0"]
test = ["0"]
```

If you have multiple IDs of the same type, separate them with a comma.

```
[mapping]
build = ["10001", "10002", "10003"]
```

You can `curl` available `issuetype`s from Jira.

```
curl \
    -D- \
    -X GET \
    -u "email@example.com:token" \
    -H "Content-Type: application/json" \
    https://{host}/rest/api/2/issuetype
```

or

```
curl \
    -D- \
    -X GET \
    -H "Authorization: Bearer {token}" \
    -H "Content-Type: application/json" \
    https://{host}/rest/api/2/issuetype
```

3. Specify the `default` which will serve as the base when checking out before deleting local branches. Specify the `origin` to be able to delete branches alongside their corresponding local branches. This ensures consistency and avoids issues during cleanup operations.

```
[branch]
default = "development"
origin  = "origin"
```

4. Specify any exclusion phrases to be removed from the branch name, if applicable.

```
[branch]
exclude = ["front","mobile","android","ios","be","web","spike","eval"]
```

5. Copy `twig/config/twig.toml` file into `~/.config/twig/` folder.

```
mkdir -p ~/.config/twig/ && \
cp twig.toml ~/.config/twig/
```

6. Compile the tool into an executable file or [download compiled executable](https://github.com/yaroslav-android/twig/releases).
> *NOTE: you might need to apply `chmod +x` to the executable if you've downloaded the precompiled version.*

```
go build
```

7. Move the executable into `/usr/local/bin` for easy global access.

```
mv twig /usr/local/bin
```

## Usage

```
twig [-h | --help] [-v | --version] [--config]
```

### twig-help

```
twig help <command>
```

Displays help information for all available commands and options in the CLI tool providing usage instructions and examples. Use this command to understand how to use the tool effectively.

#### Examples

```terminal
twig help create
```

<br/>

### twig-init

```
twig init
```

Creates a config file and folders (if they don't exist), prompts you with questions to set up the configuration interactively.

#### Examples

```terminal
twig init
```

<br/>

### twig-create

```
twig create <issue-key> [-p | --push] [-t | --type]
```

Creates the branch using Jira Issue Key as prefix after branch type.

#### Options

`-p` <br/>
`--push` - (optional) Pushes the local branch to the remote when it's created.

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

#### Examples

```terminal
twig create XX-111
```

<br/>

### twig-clean

```
twig clean local
twig clean all [-a | --assignee]
```

Deletes branches which have Jira tickets in 'Done' state.<br/>
Note: Remote branches can only be deleted if a corresponding local branch exists.

#### Options

`-a` <br/>
`--assignee` - (optional) Specifies the username (from the email) to verify that the Jira issue is assigned to you before permitting remote branch deletion. Use your Jira email, which might match `project.email`, e.g., `example.user@example.com`.

#### Examples

```
twig clean local
```
```
twig clean all -a example.user
```

<br/>

## Configuration

In case you want to experiment with custom branch formatting or extend existing methods go to `branch.go` file.

```
func BuildName(bt Type, jiraIssue network.JiraIssue, excludePhrases string) string {
    ...

    return "branchType/jiraIssue.Key_summary"
}
```

## More Examples

```
~% twig create XX-111
~% branch created: task/XX-111_jira-issue-name
```

```
~% twig create XX-111 -p
~% branch created: task/XX-111_jira-issue-name
~% remote branch created: task/XX-111_jira-issue-name
```

```
~% twig create XX-111 -t fx
~% branch created: fix/XX-111_jira-issue-name
```

```
~% twig clean local
~% branch deleted: fix/XX-111_jira-issue-name
```

```
~% twig clean all
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

```
~% twig clean all --assignee example.user
~% branch deleted: fix/XX-111_jira-issue-name
~% remote branch deleted: origin/fix/XX-111_jira-issue-name
```

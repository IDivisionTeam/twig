package common

import (
    "brcha/branch"
    "brcha/issue"
    "brcha/log"
    "brcha/network"
    "fmt"
    "strings"
)

func ConvertIssueToBranchType(issueType network.IssueType) (branch.Type, error) {
    switch issueType.Id {
    case issue.Build:
        return branch.BUILD, nil
    case issue.Chore:
        return branch.CHORE, nil
    case issue.Ci:
        return branch.CI, nil
    case issue.Docs:
        return branch.DOCS, nil
    case issue.Feat:
        return branch.FEAT, nil
    case issue.Fix:
        return branch.FIX, nil
    case issue.Perf:
        return branch.PERF, nil
    case issue.Refactor:
        return branch.REFACTOR, nil
    case issue.Revert:
        return branch.REVERT, nil
    case issue.Style:
        return branch.STYLE, nil
    case issue.Test:
        return branch.TEST, nil
    default:
        return branch.NULL, fmt.Errorf("convert: unsupported issue type %s(%s)", issueType.Name, issueType.Id)
    }
}

func ConvertUserInputToBranchType(input string) (branch.Type, error) {
    if len(input) == 0 {
        return branch.NULL, nil
    }

    switch input {
    case "build", "b":
        return branch.BUILD, nil
    case "chore", "ch":
        return branch.CHORE, nil
    case "ci":
        return branch.CI, nil
    case "docs", "d":
        return branch.DOCS, nil
    case "feat", "ft":
        return branch.FEAT, nil
    case "fix", "fx":
        return branch.FIX, nil
    case "perf", "p":
        return branch.PERF, nil
    case "refactor", "rf":
        return branch.REFACTOR, nil
    case "revert", "rv":
        return branch.REVERT, nil
    case "style", "s":
        return branch.STYLE, nil
    case "test", "t":
        return branch.TEST, nil
    default:
        return branch.NULL, fmt.Errorf("convert: unsupported branch type %s", input)
    }
}

func ConvertIssueTypesToMap(issueTypes []network.IssueType) (map[string]branch.Type, error) {
    issueMap := make(map[string]branch.Type)

    var buffer strings.Builder
    for idx, i := range issueTypes {
        _, ok := issue.Ignored[i.Id]
        if ok {
            buffer.WriteString("- ")
            buffer.WriteString(i.Name)
            buffer.WriteString("[")
            buffer.WriteString(i.Id)
            buffer.WriteString("]")
            if idx != len(issueTypes)-1 {
                buffer.WriteString("\n")
            }
            continue
        }

        name, err := ConvertIssueToBranchType(i)
        if err != nil {
            log.Warn().Printf("convert: map network to local: %v", err)
            continue
        }

        issueMap[i.Id] = name
    }

    log.Warn().Printf("convert:\nignore issue types:\n%s", buffer.String())

    return issueMap, nil
}

package common

import (
    "brcha/branch"
    "brcha/issue"
    "brcha/network"
    "fmt"
    "log"
)

func ConvertIssueToBranchType(issueType network.IssueType) (string, error) {
    branchType := branch.NewBranchType()

    var name string
    switch issueType.Id {
    case issue.Build:
        name = branchType.Build
    case issue.Chore:
        name = branchType.Chore
    case issue.Ci:
        name = branchType.Ci
    case issue.Docs:
        name = branchType.Docs
    case issue.Feat:
        name = branchType.Feat
    case issue.Fix:
        name = branchType.Fix
    case issue.Perf:
        name = branchType.Perf
    case issue.Refactor:
        name = branchType.Refactor
    case issue.Revert:
        name = branchType.Revert
    case issue.Style:
        name = branchType.Style
    case issue.Test:
        name = branchType.Test
    default:
        return "", fmt.Errorf("unsupported issue type %v", issueType)
    }

    return name, nil
}

func ConvertUserInputToBranchType(input string) (string, error) {
    branchType := branch.NewBranchType()

    if len(input) == 0 {
        return branchType.Chore, nil
    }

    var name string
    switch input {
    case "build", "b":
        name = branchType.Build
    case "chore", "ch":
        name = branchType.Chore
    case "ci":
        name = branchType.Ci
    case "docs", "d":
        name = branchType.Docs
    case "feat", "ft":
        name = branchType.Feat
    case "fix", "fx":
        name = branchType.Fix
    case "perf", "p":
        name = branchType.Perf
    case "refactor", "rf":
        name = branchType.Refactor
    case "revert", "rv":
        name = branchType.Revert
    case "style", "s":
        name = branchType.Style
    case "test", "t":
        name = branchType.Test
    default:
        return "", fmt.Errorf("unsupported branch type %s", input)
    }

    return name, nil
}

func ConvertIssueTypesToMap(issueTypes []network.IssueType) (map[string]string, error) {
    issueMap := make(map[string]string)

    for _, i := range issueTypes {
        _, ok := issue.Ignored.Get(i.Id)
        if ok {
            continue
        }

        name, err := ConvertIssueToBranchType(i)
        if err != nil {
            log.Println(err)
            continue
        }

        issueMap[i.Id] = name
    }

    return issueMap, nil
}

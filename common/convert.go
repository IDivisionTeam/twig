package common

import (
    "fmt"
    "twig/branch"
    "twig/issue"
    "twig/log"
    "twig/network"
)

func InputToBranchType(input string) (branch.Type, error) {
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
        return branch.NULL, fmt.Errorf("unsupported branch type %s", input)
    }
}

func ConvertIssueTypesToMap(issueTypes []network.IssueType) (map[string]branch.Type, error) {
    issueMap := make(map[string]branch.Type)

    local, err := issue.ParseIssueMapping()
    if err != nil {
        return nil, fmt.Errorf("convert: %w", err)
    }

    for _, i := range issueTypes {
        id, ok := local[i.Id]
        if !ok {
            log.Warn().Println(fmt.Sprintf("convert: unsupported issue type %s[%s]", i.Name, i.Id))
            continue
        }

        name, err := InputToBranchType(id)
        if err != nil {
            log.Warn().Println(fmt.Errorf("convert: %w", err))
            continue
        }

        issueMap[i.Id] = name
    }

    return issueMap, nil
}

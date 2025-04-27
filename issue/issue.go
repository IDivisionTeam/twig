package issue

import (
    "errors"
    "twig/config"
    "twig/log"
)

func ParseIssueMapping() (map[string]string, error) {
    result := make(map[string]string)
    mapping := config.GetSectionStringMap("mapping")

    if len(mapping) == 0 {
        return nil, errors.New("branch.mapping is not set")
    }

    for key, values := range mapping {
        for _, v := range values {
            if v == "0" {
                continue
            }

            result[v] = key
        }
    }

    log.Debug().Printf("issue: parsed: %v", result)
    return result, nil
}

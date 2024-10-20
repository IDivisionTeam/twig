package issue

import (
    "brcha/log"
    "strings"
)

func ParseIssueMapping(raw string) map[string]string {
    result := make(map[string]string)

    types := strings.Split(raw, ";")
    for _, t := range types {
        elements := strings.Split(t, ":")

        commitType := elements[0]
        values := strings.Split(elements[1], ",")

        for _, v := range values {
            if v == "0" {
                continue
            }
            
            result[v] = commitType
        }
    }

    log.Debug().Printf("issue: parsed: %v", result)
    return result
}

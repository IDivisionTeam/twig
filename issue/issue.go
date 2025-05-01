package issue

import (
    "fmt"
    "twig/config"
    "twig/log"
)

func ParseIssueMapping() (map[string]string, error) {
    result := make(map[string]string)
    mapping := config.GetStringMap(config.Mapping)

    if len(mapping) == 0 {
        return nil, fmt.Errorf("array %q is undefined", config.FromToken(config.Mapping))
    }

    for key, values := range mapping {
        for _, v := range values {
            if v == "0" {
                continue
            }

            result[v] = key
        }
    }

    log.Debug().Printf("mapping: %+v", result)
    return result, nil
}

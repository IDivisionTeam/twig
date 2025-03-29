package branch

import (
    "brcha/log"
    "brcha/network"
    "fmt"
    "regexp"
    "strings"
)

const (
    branchTypeSeparator string = "/"
    issueTypeSeparator  string = "_"
    wordSeparator       string = "-"
)

func BuildName(bt Type, jiraIssue network.JiraIssue, excludePhrases string) string {
    log.Info().Println("preparing branch")
    branchType := bt.ToString()
    log.Debug().Printf("build name: issue %s[%s] with branch type of '%s'", jiraIssue.Key, jiraIssue.Fields.Type.Id, branchType)

    var buffer strings.Builder

    summary := replacePhrases(*jiraIssue.Fields.Summary, excludePhrases)
    summary = camelToKebab(summary)
    summary = stripRegex(summary)

    if bt != NULL {
        buffer.WriteString(branchType)
        buffer.WriteString(branchTypeSeparator)
    }

    buffer.WriteString(jiraIssue.Key)
    buffer.WriteString(issueTypeSeparator)
    buffer.WriteString(summary)

    return buffer.String()
}

func stripRegex(in string) string {
    phrase := strings.ToLower(in)
    log.Debug().Printf("strip regex: transform: %s", phrase)

    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    phrase = reg.ReplaceAllString(phrase, wordSeparator)
    phrase = strings.TrimPrefix(phrase, wordSeparator)
    phrase = strings.TrimSuffix(phrase, wordSeparator)

    log.Debug().Printf("strip regex: transform: %s", phrase)
    return phrase
}

func replacePhrases(in string, rawPhrases string) string {
    log.Debug().Printf("replace phrases: transform: %s", in)

    phrases := strings.Split(rawPhrases, ",")

    phrase := in
    for _, v := range phrases {
        re := regexp.MustCompile("(?i)(\\[" + v + "\\]|\\(" + v + "\\))")
        phrase = re.ReplaceAllString(phrase, "")
    }

    phrase = strings.TrimSpace(phrase)
    log.Debug().Printf("replace phrases: transform: %s", phrase)
    return phrase
}

func camelToKebab(in string) string {
    log.Debug().Printf("camel2kebab: transform: %s", in)

    re1 := regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
    kebab := re1.ReplaceAllString(in, "${1}"+wordSeparator+"${2}")

    re2 := regexp.MustCompile("([a-z])([A-Z])")
    kebab = re2.ReplaceAllString(kebab, "${1}"+wordSeparator+"${2}")

    log.Debug().Printf("camel2kebab: transform: %s", kebab)
    return strings.ToLower(kebab)
}

func ExtractIssueNameFromBranch(branchName string) (string, error) {
    log.Debug().Printf("extract phrase: issue: %s", branchName)
    re := regexp.MustCompile(`[A-Z]+-\d+_`) // looking for XXXX-0000_

    match := re.FindString(branchName)
    match = strings.TrimSpace(match)
    match = strings.TrimSuffix(match, issueTypeSeparator)

    if match == "" {
        return "", fmt.Errorf("extract phrase: issue: no matches")
    }

    log.Debug().Printf("extract phrase: issue: %s", match)
    return match, nil
}

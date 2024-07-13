package branch

import (
    "brcha/log"
    "brcha/network"
    "regexp"
    "strings"
)

const (
    branchTypeSeparator string = "/"
    issueTypeSeparator  string = "_"
    wordSeparator       string = "-"
)

func BuildName(bt Type, jiraIssue network.JiraIssue) string {
    branchType := bt.ToString()
    log.Debug().Printf("build name: issue %s[%s] with branch type of %s", jiraIssue.Key, jiraIssue.Id, branchType)

    var buffer strings.Builder

    summary := replacePhrases(jiraIssue.Fields.Summary)
    summary = strings.ToLower(summary)
    summary = strings.TrimSpace(summary)
    summary = stripRegex(summary)
    summary = strings.TrimSuffix(summary, wordSeparator)

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
    log.Debug().Printf("strip regex: transform: %s", in)

    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    return reg.ReplaceAllString(in, wordSeparator)
}

func replacePhrases(in string) string {
    log.Debug().Printf("replace phrases: transform: %s", in)

    phrase := strings.ReplaceAll(in, "[Android]", "")
    phrase = strings.ReplaceAll(phrase, "[iOS]", "")
    phrase = strings.ReplaceAll(phrase, "[BE]", "")
    phrase = strings.ReplaceAll(phrase, "[WEB]", "")

    return phrase
}

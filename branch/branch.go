package branch

import (
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
    var buffer strings.Builder

    summary := replacePhrases(jiraIssue.Fields.Summary)
    summary = strings.ToLower(summary)
    summary = strings.TrimSpace(summary)
    summary = stripRegex(summary)
    summary = strings.TrimSuffix(summary, wordSeparator)

    if bt != NULL {
        buffer.WriteString(bt.ToString())
        buffer.WriteString(branchTypeSeparator)
    }

    buffer.WriteString(jiraIssue.Key)
    buffer.WriteString(issueTypeSeparator)
    buffer.WriteString(summary)

    return buffer.String()
}

func stripRegex(in string) string {
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    return reg.ReplaceAllString(in, wordSeparator)
}

func replacePhrases(in string) string {
    phrase := strings.ReplaceAll(in, "[Android]", "")
    phrase = strings.ReplaceAll(phrase, "[iOS]", "")
    phrase = strings.ReplaceAll(phrase, "[BE]", "")
    phrase = strings.ReplaceAll(phrase, "[WEB]", "")

    return phrase
}

package branch

import (
    "brcha/network"
    "regexp"
    "strings"
)

func BuildName(branchType string, jiraIssue *network.JiraIssue) string {
    var buffer strings.Builder

    summary := replacePhrases(jiraIssue.Fields.Summary)
    summary = strings.ToLower(summary)
    summary = strings.TrimSpace(summary)
    summary = stripRegex(summary)
    summary = strings.TrimSuffix(summary, "-")

    if len(branchType) > 0 {
        buffer.WriteString(branchType)
        buffer.WriteString("/")
    }

    buffer.WriteString(jiraIssue.Key)
    buffer.WriteString("_")
    buffer.WriteString(summary)

    return buffer.String()
}

func stripRegex(in string) string {
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    return reg.ReplaceAllString(in, "-")
}

func replacePhrases(in string) string {
    return strings.ReplaceAll(in, "[Android]", "")
}

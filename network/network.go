package network

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

type jiraCredentials struct {
    host  string
    email string
    token string
}

// TODO: replace with your <host>, <email>, and <token>.
func withJiraCredentials() *jiraCredentials {
    return &jiraCredentials{
        host:  "",
        email: "",
        token: "",
    }
}

type JiraIssue struct {
    Id     string      `json:"id"`
    Key    string      `json:"key"`
    Fields IssueFields `json:"fields"`
}

type IssueFields struct {
    Type    IssueType `json:"issuetype"`
    Summary string    `json:"summary"`
}

type IssueType struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

func GetJiraIssue(issueKey string) (*JiraIssue, error) {
    jira := withJiraCredentials()
    url := fmt.Sprintf("https://%s/rest/api/2/issue/%s?fields=issuetype,summary", jira.host, issueKey)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    request.SetBasicAuth(jira.email, jira.token)
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Print(err)
        }
    }(response.Body)

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("jira API not available %w", err)
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    var jiraIssue JiraIssue
    err = json.Unmarshal(body, &jiraIssue)
    if err != nil {
        return nil, err
    }

    return &jiraIssue, nil
}

func GetJiraIssueTypes() ([]IssueType, error) {
    jira := withJiraCredentials()
    url := fmt.Sprintf("https://%s/rest/api/2/issuetype", jira.host)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    request.SetBasicAuth(jira.email, jira.token)
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Print(err)
        }
    }(response.Body)

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("jira API not available %w", err)
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    var jiraIssue []IssueType
    err = json.Unmarshal(body, &jiraIssue)
    if err != nil {
        return nil, err
    }

    return jiraIssue, nil
}

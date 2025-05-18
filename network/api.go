package network

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "twig/log"
)

type JiraApi interface {
    GetJiraIssueTypes() ([]IssueType, error)
    GetJiraIssue(issueKey string) (*JiraIssue, error)
    GetJiraIssueStatus(issueKey string, hasAssignee bool) (*JiraIssue, error)
    GetJiraIssueStatusBulk(issueKeys []string, hasAssignee bool) ([]JiraIssue, error)
}

type mixedJiraApi struct {
    client Client
}

func NewJiraApi(client Client) JiraApi {
    return &mixedJiraApi{
        client: client,
    }
}

func (api *mixedJiraApi) GetJiraIssueTypes() ([]IssueType, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issuetype'", http.MethodGet))
    path := "issuetype"

    response, err := api.client.SendRequest(http.MethodGet, path, nil)
    if err != nil {
        return nil, err
    }

    log.Debug().Println(fmt.Sprintf("Response %d 'issuetype'\n%s", response.statusCode, response.body))

    var jiraIssue []IssueType
    err = json.Unmarshal(response.body, &jiraIssue)
    if err != nil {
        return nil, err
    }

    return jiraIssue, nil
}

func (api *mixedJiraApi) GetJiraIssue(issueKey string) (*JiraIssue, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issue'", http.MethodGet))
    path := fmt.Sprintf("issue/%s?fields=issuetype,summary", issueKey)

    response, err := api.client.SendRequest(http.MethodGet, path, nil)
    if err != nil {
        return nil, err
    }

    log.Debug().Println(fmt.Sprintf("Response %d 'issue'\n%s", response.statusCode, response.body))

    var jiraIssue JiraIssue
    if err := json.Unmarshal(response.body, &jiraIssue); err != nil {
        return nil, err
    }

    return &jiraIssue, nil
}

func (api *mixedJiraApi) GetJiraIssueStatus(issueKey string, hasAssignee bool) (*JiraIssue, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issue status'", http.MethodGet))

    path := fmt.Sprintf("issue/%s?fields=status", issueKey)
    if hasAssignee {
        path = fmt.Sprintf("%s%s", path, ",assignee")
    }

    response, err := api.client.SendRequest(http.MethodGet, path, nil)
    if err != nil {
        return nil, err
    }

    log.Debug().Println(fmt.Sprintf("Response %d 'issue status'\n%s", response.statusCode, response.body))

    var jiraIssue JiraIssue
    if err := json.Unmarshal(response.body, &jiraIssue); err != nil {
        return nil, err
    }

    return &jiraIssue, nil
}

func (api *mixedJiraApi) GetJiraIssueStatusBulk(issueKeys []string, hasAssignee bool) ([]JiraIssue, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issue status bulk'", http.MethodPost))

    fields := []string{"status"}
    if hasAssignee {
        fields = append(fields, "assignee")
    }

    body := JiraIssueBulkRequest{
        Fields:    fields,
        IssueKeys: issueKeys,
    }

    log.Debug().Printf(fmt.Sprintf("Request body\n%+v", body))

    encodedBody, _ := json.Marshal(body)
    response, err := api.client.SendRequest(http.MethodPost, "issue/bulkfetch", bytes.NewBuffer(encodedBody))
    if err != nil {
        return nil, err
    }

    log.Debug().Println(fmt.Sprintf("Response %d 'issue status bulk'\n%s", response.statusCode, response.body))

    var jiraIssues JiraIssues
    if err := json.Unmarshal(response.body, &jiraIssues); err != nil {
        return nil, err
    }

    return jiraIssues.Issues, nil
}

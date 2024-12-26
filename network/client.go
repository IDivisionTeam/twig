package network

import (
    "brcha/log"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
)

type jiraCredentials struct {
    host  string
    email string
    token string
}

func readJiraCredentials() *jiraCredentials {
    return &jiraCredentials{
        host:  os.Getenv("BRCHA_HOST"),
        email: os.Getenv("BRCHA_EMAIL"),
        token: os.Getenv("BRCHA_TOKEN"),
    }
}

type Client interface {
    GetJiraIssueTypes() ([]IssueType, error)
    GetJiraIssue(issueKey string) (*JiraIssue, error)
    GetJiraIssueStatus(issueKey string) (*JiraIssue, error)
}

type networkClient struct {
    credentials *jiraCredentials
    client      *http.Client
}

func NewClient(client *http.Client) Client {
    return &networkClient{
        credentials: readJiraCredentials(),
        client:      client,
    }
}

func (c *networkClient) GetJiraIssueTypes() ([]IssueType, error) {
    log.Info().Println("sending request <issuetype>")
    path := "issuetype"

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue types: %w", err)
    }

    log.Debug().Printf("response <issuetype>:\n%s", response.body)

    var jiraIssue []IssueType
    err = json.Unmarshal(response.body, &jiraIssue)
    if err != nil {
        return nil, fmt.Errorf("get issue: (%d) types unmarshal: %w", response.statusCode, err)
    }

    return jiraIssue, nil
}

func (c *networkClient) GetJiraIssue(issueKey string) (*JiraIssue, error) {
    log.Info().Println("sending request <issue>")
    path := fmt.Sprintf("issue/%s?fields=issuetype,summary", issueKey)

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue: %w", err)
    }

    log.Debug().Printf("response <issue>:\n%s", response.body)

    var jiraIssue JiraIssue
    if err := json.Unmarshal(response.body, &jiraIssue); err != nil {
        return nil, fmt.Errorf("get issue: (%d) unmarshal: %w", response.statusCode, err)
    }

    return &jiraIssue, nil
}

func (c *networkClient) GetJiraIssueStatus(issueKey string) (*JiraIssue, error) {
    log.Info().Println("sending request <issue-status>")
    path := fmt.Sprintf("issue/%s?fields=status", issueKey)

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue: %w", err)
    }

    log.Debug().Printf("response <issue-status>:\n%s", response.body)

    var jiraIssue JiraIssue
    if err := json.Unmarshal(response.body, &jiraIssue); err != nil {
        return nil, fmt.Errorf("get issue: (%d) unmarshal: %w", response.statusCode, err)
    }

    return &jiraIssue, nil
}

func (c *networkClient) prepareRequest(path string) (*http.Request, error) {
    url := fmt.Sprintf("https://%s/rest/api/2/%s", c.credentials.host, path)
    log.Debug().Printf("prepare request: url = %s", url)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("prepare request: %w", err)
    }

    addAuthHeader(request, c.credentials)
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    return request, nil
}

func addAuthHeader(request *http.Request, credentials *jiraCredentials) {
    isBasicAuth := len(credentials.email) > 0

    if isBasicAuth {
        log.Debug().Println("prepare request: basic auth")
        request.SetBasicAuth(credentials.email, credentials.token)
        return
    }

    log.Debug().Println("prepare request: bearer auth")
    bearer := "Bearer " + credentials.token
    request.Header.Add("Authorization", bearer)
}

func (c *networkClient) sendRequest(path string) (*Response, error) {
    request, err := c.prepareRequest(path)
    if err != nil {
        return nil, fmt.Errorf("send request: %w", err)
    }

    log.Debug().Println("send request: enqueue request")
    response, err := c.client.Do(request)
    if err != nil {
        log.Warn().Println("verify auth type (basic or bearer)")
        return nil, fmt.Errorf("send request: %w", err)
    }

    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Error().Printf("send request: %v", err)
        }
    }(response.Body)

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("send request: read response bytes: %w", err)
    }

    log.Info().Printf("status code = %d", response.StatusCode)
    if response.StatusCode == http.StatusOK {
        return &Response{
            statusCode: response.StatusCode,
            body:       body,
        }, nil
    }

    var jiraError JiraError
    if err := json.Unmarshal(body, &jiraError); err != nil {
        return nil, fmt.Errorf("send request: unmarshal error: %w", err)
    }

    errors := strings.Join(jiraError.ErrorMessages[:], "\n")
    return nil, fmt.Errorf("send request: Jira API: %s", errors)
}

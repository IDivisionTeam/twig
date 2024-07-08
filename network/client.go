package network

import (
    "brcha/recorder"
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

func withJiraCredentials() *jiraCredentials {
    return &jiraCredentials{
        host:  os.Getenv("BRCHA_HOST"),
        email: os.Getenv("BRCHA_EMAIL"),
        token: os.Getenv("BRCHA_TOKEN"),
    }
}

type Client interface {
    GetJiraIssueTypes() ([]IssueType, error)
    GetJiraIssue(issueKey string) (*JiraIssue, error)
}

type networkClient struct {
    credentials *jiraCredentials
    client      *http.Client
}

func NewClient(client *http.Client) Client {
    return &networkClient{
        credentials: withJiraCredentials(),
        client:      client,
    }
}

func (c *networkClient) GetJiraIssueTypes() ([]IssueType, error) {
    path := "issuetype"

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue types: %w", err)
    }

    var jiraIssue []IssueType
    err = json.Unmarshal(response.body, &jiraIssue)
    if err != nil {
        return nil, fmt.Errorf("(%d) get issue types unmarshal: %w", response.statusCode, err)
    }

    return jiraIssue, nil
}

func (c *networkClient) GetJiraIssue(issueKey string) (*JiraIssue, error) {
    path := fmt.Sprintf("issue/%s?fields=issuetype,summary", issueKey)

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue: %w", err)
    }

    var jiraIssue JiraIssue
    if err := json.Unmarshal(response.body, &jiraIssue); err != nil {
        return nil, fmt.Errorf("(%d) get issue unmarshal: %w", response.statusCode, err)
    }

    return &jiraIssue, nil
}

func (c *networkClient) prepareRequest(path string) (*http.Request, error) {
    url := fmt.Sprintf("https://%s/rest/api/2/%s", c.credentials.host, path)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("prepare request: %w", err)
    }

    request.SetBasicAuth(c.credentials.email, c.credentials.token)
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    return request, nil
}

func (c *networkClient) sendRequest(path string) (*Response, error) {
    request, err := c.prepareRequest(path)
    if err != nil {
        return nil, fmt.Errorf("new request: %w", err)
    }

    response, err := c.client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("send request: %w", err)
    }

    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            recorder.Println(recorder.ERROR, err)
        }
    }(response.Body)

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("read response bytes: %w", err)
    }

    if response.StatusCode != http.StatusOK {
        var jiraError JiraError
        if err := json.Unmarshal(body, &jiraError); err != nil {
            return nil, fmt.Errorf("(%d) unmarshal error: %w", response.StatusCode, err)
        }

        errors := strings.Join(jiraError.ErrorMessages[:], "\n")

        return nil, fmt.Errorf("(%d) Jira API: %s", response.StatusCode, errors)
    }

    return &Response{
        statusCode: response.StatusCode,
        body:       body,
    }, nil
}

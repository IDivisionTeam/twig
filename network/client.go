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
    log.Debug().Println("get issuetype: request issue types")
    path := "issuetype"

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue types: %w", err)
    }

    var jiraIssue []IssueType
    err = json.Unmarshal(response.body, &jiraIssue)
    if err != nil {
        return nil, fmt.Errorf("get issue: (%d) types unmarshal: %w", response.statusCode, err)
    }

    return jiraIssue, nil
}

func (c *networkClient) GetJiraIssue(issueKey string) (*JiraIssue, error) {
    log.Debug().Println("get issue: request issue")
    path := fmt.Sprintf("issue/%s?fields=issuetype,summary", issueKey)

    response, err := c.sendRequest(path)
    if err != nil {
        return nil, fmt.Errorf("get issue: %w", err)
    }

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
        request.SetBasicAuth(credentials.email, credentials.token)
        return
    }

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

    log.Debug().Printf("send request: status code = %d", response.StatusCode)
    if response.StatusCode != http.StatusOK {
        var jiraError JiraError
        if err := json.Unmarshal(body, &jiraError); err != nil {
            return nil, fmt.Errorf("send request: (%d) unmarshal error: %w", response.StatusCode, err)
        }

        errors := strings.Join(jiraError.ErrorMessages[:], "\n")

        return nil, fmt.Errorf("send request: (%d) Jira API: %s", response.StatusCode, errors)
    }

    return &Response{
        statusCode: response.StatusCode,
        body:       body,
    }, nil
}

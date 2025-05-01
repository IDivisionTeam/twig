package network

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "strings"
    "twig/config"
    "twig/log"
)

type jiraCredentials struct {
    host  string
    auth  string
    email string
    token string
}

func readJiraCredentials() *jiraCredentials {
    return &jiraCredentials{
        host:  config.GetString(config.ProjectHost),
        auth:  config.GetString(config.ProjectAuth),
        email: config.GetString(config.ProjectEmail),
        token: config.GetString(config.ProjectToken),
    }
}

type Client interface {
    GetJiraIssueTypes() ([]IssueType, error)
    GetJiraIssue(issueKey string) (*JiraIssue, error)
    GetJiraIssueStatus(issueKey string, hasAssignee bool) (*JiraIssue, error)
    GetJiraIssueStatusBulk(issueKeys []string, hasAssignee bool) ([]JiraIssue, error)
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
    log.Debug().Println(fmt.Sprintf("Request %s 'issuetype'", http.MethodGet))
    path := "issuetype"

    response, err := c.sendRequest(http.MethodGet, path, nil)
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

func (c *networkClient) GetJiraIssue(issueKey string) (*JiraIssue, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issue'", http.MethodGet))
    path := fmt.Sprintf("issue/%s?fields=issuetype,summary", issueKey)

    response, err := c.sendRequest(http.MethodGet, path, nil)
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

func (c *networkClient) GetJiraIssueStatus(issueKey string, hasAssignee bool) (*JiraIssue, error) {
    log.Debug().Println(fmt.Sprintf("Request %s 'issue status'", http.MethodGet))

    path := fmt.Sprintf("issue/%s?fields=status", issueKey)
    if hasAssignee {
        path = fmt.Sprintf("%s%s", path, ",assignee")
    }

    response, err := c.sendRequest(http.MethodGet, path, nil)
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

func (c *networkClient) GetJiraIssueStatusBulk(issueKeys []string, hasAssignee bool) ([]JiraIssue, error) {
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
    response, err := c.sendRequest(http.MethodPost, "issue/bulkfetch", bytes.NewBuffer(encodedBody))
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

func (c *networkClient) prepareRequest(method, path string, body io.Reader) (*http.Request, error) {
    url := fmt.Sprintf("https://%s/rest/api/2/%s", c.credentials.host, path)
    log.Debug().Println(fmt.Sprintf("Request path %q", path))

    request, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    if err := addAuthHeader(request, c.credentials); err != nil {
        return nil, err
    }

    if method == http.MethodPost {
        request.Header.Set("Accept", "application/json")
    }
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    return request, nil
}

func addAuthHeader(request *http.Request, credentials *jiraCredentials) error {
    auth := strings.ToLower(credentials.auth)

    switch auth {
    case basicType:
        log.Debug().Println("Use Basic Auth")
        request.SetBasicAuth(credentials.email, credentials.token)
        return nil
    case bearerType:
        log.Debug().Println("Use Bearer Auth")
        request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", credentials.token))
        return nil
    default:
        return fmt.Errorf("%q does not support %q auth", config.FromToken(config.ProjectAuth), auth)
    }
}

func (c *networkClient) sendRequest(method, path string, body io.Reader) (*Response, error) {
    request, err := c.prepareRequest(method, path, body)
    if err != nil {
        return nil, err
    }

    log.Debug().Println(fmt.Sprintf("Enqueue request %q", path))
    response, err := c.client.Do(request)
    if err != nil {
        log.Warn().Println("Verify auth type (Basic/Bearer)")
        return nil, err
    }

    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            log.Error().Println(fmt.Errorf("request: %w", err))
        }
    }(response.Body)

    data, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    if response.StatusCode == http.StatusOK {
        return &Response{
            statusCode: response.StatusCode,
            body:       data,
        }, nil
    }

    var jiraError JiraError
    if err := json.Unmarshal(data, &jiraError); err != nil {
        return nil, err
    }

    errs := strings.Join(jiraError.ErrorMessages[:], "\n")
    return nil, errors.New(errs)
}

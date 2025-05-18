package network

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "strings"
    "twig/config"
    "twig/log"
)

type Client interface {
    PrepareRequest(method, path string, body io.Reader) (*http.Request, error)
    SendRequest(method, path string, body io.Reader) (*Response, error)
}

type httpClient struct {
    credentials *jiraCredentials
    client      *http.Client
}

func NewHttpClient(client *http.Client) Client {
    return &httpClient{
        credentials: &jiraCredentials{
            host:  config.GetString(config.ProjectHost),
            auth:  config.GetString(config.ProjectAuth),
            email: config.GetString(config.ProjectEmail),
            token: config.GetString(config.ProjectToken),
        },
        client: client,
    }
}

func (c *httpClient) PrepareRequest(method, path string, body io.Reader) (*http.Request, error) {
    url := fmt.Sprintf("https://%s/rest/api/2/%s", c.credentials.host, path)
    log.Debug().Println(fmt.Sprintf("Request path %q", path))

    request, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    if err = c.cofigureHeaders(method, request); err != nil {
        return nil, err
    }

    return request, nil
}

func (c *httpClient) SendRequest(method, path string, body io.Reader) (*Response, error) {
    request, err := c.PrepareRequest(method, path, body)
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

func (c *httpClient) cofigureHeaders(method string, request *http.Request) error {
    if err := c.addAuthHeader(request, c.credentials); err != nil {
        return err
    }

    if method == http.MethodPost {
        request.Header.Set("Accept", "application/json")
    }
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")

    return nil
}

func (c *httpClient) addAuthHeader(request *http.Request, credentials *jiraCredentials) error {
    auth := strings.ToLower(credentials.auth)

    switch auth {
    case BasicType:
        log.Debug().Println("Use Basic Auth")
        request.SetBasicAuth(credentials.email, credentials.token)
        return nil
    case BearerType:
        log.Debug().Println("Use Bearer Auth")
        request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", credentials.token))
        return nil
    default:
        return fmt.Errorf("%q does not support %q auth", config.FromToken(config.ProjectAuth), auth)
    }
}

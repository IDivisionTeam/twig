package network

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

type JiraError struct {
    ErrorMessages []string `json:"errorMessages"`
}

type Response struct {
    statusCode int
    body       []byte
}

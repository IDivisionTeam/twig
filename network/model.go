package network

const (
    BasicType  = "basic"
    BearerType = "bearer"
)

type jiraCredentials struct {
    host  string
    auth  string
    email string
    token string
}

type JiraIssues struct {
    Issues []JiraIssue `json:"issues"`
}

type JiraIssue struct {
    Id     string      `json:"id"`
    Key    string      `json:"key"`
    Fields IssueFields `json:"fields"`
}

type IssueFields struct {
    Type     *IssueType     `json:"issuetype,omitempty"`
    Summary  *string        `json:"summary,omitempty"`
    Status   *IssueStatus   `json:"status,omitempty"`
    Assignee *IssueAssignee `json:"assignee,omitempty"`
}

type IssueType struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

type IssueStatus struct {
    Category IssueStatusCategory `json:"statusCategory"`
}

type IssueStatusCategory struct {
    Id   int    `json:"id"`
    Name string `json:"key"`
}

type IssueAssignee struct {
    Email string `json:"emailAddress"`
}

type JiraError struct {
    ErrorMessages []string `json:"errorMessages"`
}

type Response struct {
    statusCode int
    body       []byte
}

type JiraIssueBulkRequest struct {
    Fields    []string `json:"fields"`
    IssueKeys []string `json:"issueIdsOrKeys"`
}

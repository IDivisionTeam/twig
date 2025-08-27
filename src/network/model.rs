use serde::{Deserialize, Serialize};

pub const AUTH_BASIC: &str = "basic";
pub const AUTH_BEARER: &str = "bearer";

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct JiraIssues {
    issues: Vec<JiraIssue>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct JiraIssue {
    id: String,
    key: String,
    fields: JiraIssueFields,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct JiraIssueFields {
    #[serde(rename = "issuetype", skip_serializing_if = "Option::is_none")]
    issue_type: Option<JiraIssueType>,

    #[serde(skip_serializing_if = "Option::is_none")]
    summary: Option<String>,

    #[serde(skip_serializing_if = "Option::is_none")]
    status: Option<JiraIssueStatus>,

    #[serde(skip_serializing_if = "Option::is_none")]
    assignee: Option<JiraIssueAssignee>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraIssueType {
    id: String,
    name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraIssueStatus {
    #[serde(rename = "statusCategory")]
    category: JiraIssueStatusCategory,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraIssueStatusCategory {
    id: String,

    #[serde(rename = "key")]
    name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraIssueAssignee {
    #[serde(rename = "emailAddress")]
    email: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraIssueBulkRequest {
    fields: Vec<String>,

    #[serde(rename = "issueIdsOrKeys")]
    issue_keys: Vec<String>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
struct JiraError {
    #[serde(rename = "errorMessages")]
    messages: Option<Vec<String>>,
}

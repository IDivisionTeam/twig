use serde::{Deserialize, Serialize};

pub const AUTH_BASIC: &str = "basic";
pub const AUTH_BEARER: &str = "bearer";

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssues {
    issues: Vec<JiraIssue>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssue {
    id: String,
    key: String,
    fields: JiraIssueFields,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssueFields {
    #[serde(rename = "issuetype", skip_serializing_if = "Option::is_none")]
    issue_type: Option<JiraIssueType>,

    #[serde(skip_serializing_if = "Option::is_none")] // Для чого skip_serializing_if?
    summary: Option<String>,

    status: Option<JiraIssueStatus>,

    #[serde(skip_serializing_if = "Option::is_none")]
    assignee: Option<JiraIssueAssignee>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueType {
    id: String,
    name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatus {
    #[serde(rename = "statusCategory")]
    category: JiraIssueStatusCategory,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusFields {
    status: JiraIssueStatus,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusObject {
    fields: JiraIssueStatusFields,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusBulkResponse {
    issues: Vec<JiraIssueStatusObject>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusCategory {
    id: i64,

    #[serde(rename = "key")]
    name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueAssignee {
    #[serde(rename = "emailAddress")]
    email: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueBulkRequest {
    pub fields: Vec<String>,

    #[serde(rename = "issueIdsOrKeys")]
    pub issue_keys: Vec<String>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraError {
    #[serde(rename = "errorMessages")]
    messages: Option<Vec<String>>,
}

pub struct Credentials {
    pub host: String,
    pub auth: String,
    pub email: String,
    pub token: String,
}

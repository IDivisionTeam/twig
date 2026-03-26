#![allow(dead_code)]

use serde::{Deserialize, Serialize};

use crate::network::client::ApiError;

pub const AUTH_BASIC: &str = "basic";
pub const AUTH_BEARER: &str = "bearer";

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssues {
    pub issues: Vec<JiraIssue>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssue {
    pub id: String,
    pub key: String,
    pub fields: JiraIssueFields,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JiraIssueFields {
    #[serde(rename = "issuetype", skip_serializing_if = "Option::is_none")]
    pub issue_type: Option<JiraIssueType>,

    #[serde(skip_serializing_if = "Option::is_none")] // Для чого skip_serializing_if?
    pub summary: Option<String>,

    pub status: Option<JiraIssueStatus>,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub assignee: Option<JiraIssueAssignee>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueType {
    pub id: String,
    pub name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatus {
    #[serde(rename = "statusCategory")]
    pub category: JiraIssueStatusCategory,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusFields {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub assignee: Option<JiraIssueAssignee>,
    pub status: JiraIssueStatus,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusObject {
    pub key: String,
    pub fields: JiraIssueStatusFields,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusBulkResponse {
    pub issues: Vec<JiraIssueStatusObject>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueStatusCategory {
    pub id: i64,

    #[serde(rename = "key")]
    pub name: String,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct JiraIssueAssignee {
    #[serde(rename = "emailAddress")]
    pub email: String,
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
    pub messages: Option<Vec<String>>,
}

impl ApiError for JiraError {
    fn errors(&self) -> Vec<String> {
        self.messages.clone().unwrap_or_default()
    }
}

pub struct Credentials {
    pub host: String,
    pub auth: String,
    pub email: String,
    pub token: String,
}

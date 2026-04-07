use anyhow::{Context, Result};
use reqwest::Method;
use serde::{Serialize, de::DeserializeOwned};

use crate::network::{
    client::TwigClient,
    model::{
        Credentials, JiraError, JiraIssue, JiraIssueBulkRequest, JiraIssueStatusBulkResponse,
        JiraIssueStatusObject, JiraIssueType,
    },
};
use crate::network::client::ApiError;

pub trait HttpClient {
    fn get<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
        &self,
        path: &str,
        params: Vec<(&str, &str)>,
    ) -> Result<T>;

    fn post<T: DeserializeOwned, S: Serialize, E: ApiError + DeserializeOwned>(
        &self,
        path: &str,
        body: S,
    ) -> Result<T>;

}
pub struct JiraApi<C: HttpClient> {
    client: C,
}

impl<C: HttpClient> JiraApi<C> {
    pub fn new(client: C) -> JiraApi<C> {
        Self { client }
    }

    fn get<T: DeserializeOwned>(&self, path: &str, params: Vec<(&str, &str)>) -> Result<T> {
        self.client.get::<T, JiraError>(path, params)
    }

    fn post<T: DeserializeOwned, S: Serialize>(&self, path: &str, body: S) -> Result<T> {
        self.client.post::<T, S, JiraError>(path, body)
    }

    #[allow(dead_code)]
    pub fn get_jira_issue_types(&self) -> Result<Vec<JiraIssueType>> {
        self.get("issuetype", vec![])
            .context("failed to get jira issues types")
    }

    #[allow(dead_code)]
    pub fn get_jira_issue(&self, issue_key: &str) -> Result<JiraIssue> {
        self.get(
            &format!("issue/{issue_key}"),
            vec![("fields", "issuetype,summary")],
        )
        .context("failed to get jira issues types")
    }

    #[allow(dead_code)]
    pub fn get_jira_issues(&self, _issue_keys: Vec<&str>) -> Result<Vec<JiraIssue>> {
        todo!()
    }

    pub fn get_jira_issue_status(
        &self,
        issue_key: &str,
        has_assignee: bool,
    ) -> Result<JiraIssueStatusObject> {
        let mut fields = "status".to_string();
        if has_assignee {
            fields = format!("{fields},assignee");
        }
        self.get(&format!("issue/{issue_key}"), vec![("fields", &fields)])
            .context("failed to get jira issues types")
    }

    #[allow(dead_code)]
    pub fn get_jira_issue_statuses(
        &self,
        issue_keys: Vec<String>,
        has_assignee: bool,
    ) -> Result<JiraIssueStatusBulkResponse> {
        let mut fields = vec!["status".to_string()];
        if has_assignee {
            fields.push("assignee".to_string());
        }

        let body = JiraIssueBulkRequest { fields, issue_keys };
        self.post("issue/bulkfetch", body)
            .context("failed to get jira issues statutes")
    }
}

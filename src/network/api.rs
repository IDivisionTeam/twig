use anyhow::{Context, Result};

use crate::network::{
    client::TwigClient,
    model::{
        Credentials, JiraIssue, JiraIssueBulkRequest, JiraIssueStatusBulkResponse, JiraIssueStatusObject, JiraIssueType
    },
};

pub struct JiraApi {
    client: TwigClient,
}

impl JiraApi {
    pub fn new(credentials: Credentials) -> Result<JiraApi> {
        let client = TwigClient::new(credentials)?;
        Ok(Self { client })
    }

    pub fn get_jira_issue_types(&self) -> Result<Vec<JiraIssueType>> {
        self.client
            .get("issuetype", vec![])
            .context("failed to get jira issues types")
    }

    pub fn get_jira_issue(&self, issue_key: &str) -> Result<JiraIssue> {
        self.client
            .get(
                &format!("issue/{issue_key}"),
                vec![("fields", "issuetype,summary")],
            )
            .context("failed to get jira issues types")
    }

    pub fn get_jira_issues(&self, issue_keys: Vec<&str>) -> Result<Vec<JiraIssue>> { // У го такого метода немає 🤷
        todo!()
    }

    pub fn get_jira_issue_status(
        &self,
        issue_key: String,
        has_assignee: bool,
    ) -> Result<JiraIssueStatusObject> {
        let mut fields = "status".to_string();
        if has_assignee {
            fields = format!("{fields},assignee");
        }
        self.client
            .get(&format!("issue/{issue_key}"), vec![("fields", &fields)])
            .context("failed to get jira issues types")
    }

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
        self.client
            .post("issue/bulkfetch", body)
            .context("failed to get jira issues statutes")
    }
}

use crate::network::model::{JiraIssue, JiraIssueStatus, JiraIssueType};

pub struct JiraApi {}

impl JiraApi {
    pub fn new() -> JiraApi {
        Self {}
    }

    pub fn get_jira_issue_types() -> Result<Vec<JiraIssueType>, &'static str> {
        Err("not implemented")
    }

    pub fn get_jira_issue(issue_key: String) -> Result<JiraIssue, &'static str> {
        Err("not implemented")
    }

    pub fn get_jira_issues(issue_keys: Vec<String>) -> Result<Vec<JiraIssue>, &'static str> {
        Err("not implemented")
    }

    pub fn get_jira_issue_status(
        issue_key: String,
        has_assignee: bool,
    ) -> Result<JiraIssueStatus, &'static str> {
        Err("not implemented")
    }

    pub fn get_jira_issue_statuses(
        issue_keys: Vec<String>,
        has_assignee: bool,
    ) -> Result<Vec<JiraIssueStatus>, &'static str> {
        Err("not implemented")
    }
}

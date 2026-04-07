use std::{collections::HashMap, sync::LazyLock};

use clap::{Args, Subcommand};
use log::{debug, info};
use regex::Regex;

use crate::{
    branch, config, git,
    network::{self, api::JiraApi},
};

const ITEMS_THRESHOLD: usize = 5;
const DONE_STATUS_ID: i64 = 3;

static ISSUE_REGEX: LazyLock<Regex> = LazyLock::new(|| Regex::new(r"[A-Z]+-\d+_").unwrap());

use thiserror::Error;
use crate::network::api::HttpClient;

#[derive(Error, Debug)]
pub enum CleanCmdError {
    #[error("{}", Self::no_done_issue_found(.0))]
    NoDoneIssueFound(Option<String>),
    #[error("email {0} is either invalid or corrupted")]
    InvalidEmail(String),
    #[error("no branches related to Jira issues were found")]
    NoIssuesForBranches,
    #[error("nothing to clean")]
    NothingToClean,
    #[error("failed to extract issue fromb branch '{0}'")]
    ExtractIssueFromBranch(String),
    #[error("validate: issue '{issue_key}' has assignee '{username}' but looking for '{assignee}'")]
    InvalidAssignee {
        issue_key: String,
        username: String,
        assignee: String,
    },
    #[error("git failed")]
    GitExecute(#[from] git::GitError),
}

impl CleanCmdError {
    fn no_done_issue_found(assignee: &Option<String>) -> String {
        match assignee {
            Some(assignee) => {
                format!("no associated Jira issues in DONE status where assignee is '{assignee}'")
            }
            None => "no associated Jira issues in DONE status".to_string(),
        }
    }
}

type Result<T> = std::result::Result<T, CleanCmdError>;

#[derive(Args)]
pub struct Clean {
    /// "(optional) overrides the assignee used when comparing before deleting the branch, default is username from credentials.email"
    #[arg(short, long)]
    assignee: Option<String>,

    /// (optional) delete branch while ignoring the assignee; the 'assignee' option is disregarded when this flag is used
    #[arg(short, long, default_value_t = false)]
    ignore_assignee: bool,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Deletes only local branches which have Jira tickets in 'Done' state.
    Local,
    /// Deletes remote and local branches which have Jira tickets in 'Done' state.
    All,
}

pub fn handle<C: HttpClient>(jira_api: &JiraApi<C>, args: &Clean, config: &config::Config) -> Result<()> {
    let delete_remote = match args.command {
        Commands::Local => false,
        Commands::All => true,
    };

    let user_name = extract_username_from_email(&config.credentials.email)?;
    let assignee = args.assignee.as_ref().unwrap_or(&user_name);

    let fetch_output = git::fetch_prune()?;
    if !fetch_output.is_empty() {
        info!("{fetch_output}");
    }
    git::ensure_clean_status()?;
    let checkout_output = git::checkout(&config.project.branch)?;
    if !checkout_output.is_empty() {
        info!("{checkout_output}");
    }

    let local_branches = git::get_local_branches()?;
    let issues = pair_branches_with_issues(&local_branches)?;
    let statuses = pair_branches_with_statuses(jira_api, issues, assignee, args.ignore_assignee)?;
    let deleted_any = delete_branches_if_any(delete_remote, &config.project.remote, statuses);

    if !deleted_any {
        return Err(CleanCmdError::NoDoneIssueFound(
            (!args.ignore_assignee).then_some(assignee.clone()),
        ));
    }

    Ok(())
}

fn extract_username_from_email(email: &str) -> Result<String> {
    let parts = email.split("@").collect::<Vec<_>>();
    if parts.len() != 2 {
        return Err(CleanCmdError::InvalidEmail(email.to_string()));
    }

    Ok(parts[0].to_string())
}

fn pair_branches_with_issues(raw_branches: &str) -> Result<HashMap<String, String>> {
    let mut issues = HashMap::new();

    for local_branch in raw_branches.split_whitespace() {
        let trimmed_branch_name = local_branch.split_whitespace().collect::<Vec<_>>().join("");
        let Some(issue) = extract_issue_name_from_branch(&trimmed_branch_name)
            .ok()
            .filter(|s| !s.is_empty())
        else {
            continue;
        };

        debug!("Branch '{issue}' with issue '{trimmed_branch_name}'");

        issues.insert(trimmed_branch_name, issue);
    }

    if issues.is_empty() {
        return Err(CleanCmdError::NoIssuesForBranches);
    }

    Ok(issues)
}

fn pair_branches_with_statuses<C: HttpClient>(
    jira_api: &JiraApi<C>,
    issues: HashMap<String, String>,
    assignee: &str,
    ignore_assignee: bool,
) -> Result<HashMap<String, network::model::JiraIssueStatusCategory>> {
    let statuses = if issues.len() < ITEMS_THRESHOLD {
        query_issues(jira_api, issues, assignee, ignore_assignee)
    } else {
        // TODO: [BR-61] implement bulk_query_issues and call it instead of query_issues
        query_issues(jira_api, issues, assignee, ignore_assignee)
    };

    if statuses.is_empty() {
        return Err(CleanCmdError::NothingToClean);
    }

    Ok(statuses)
}

fn delete_branches_if_any(
    delete_remote: bool,
    remote: &str,
    statuses: HashMap<String, network::model::JiraIssueStatusCategory>,
) -> bool {
    let mut any_in_done_status = false;

    for (branch_name, status) in statuses {
        if status.id == DONE_STATUS_ID {
            delete_local_branch(&branch_name);

            if delete_remote {
                delete_remote_branch(remote, &branch_name);
            }

            any_in_done_status = true;
        }
    }

    any_in_done_status
}

fn query_issues<C: HttpClient>(
    jira_api: &JiraApi<C>,
    issues: HashMap<String, String>,
    assignee: &str,
    ignore_assignee: bool,
) -> HashMap<String, network::model::JiraIssueStatusCategory> {
    let mut statuses = HashMap::new();

    for (local_branch, issue) in issues {
        let Some(jira_issue) = jira_api
            .get_jira_issue_status(&issue, !ignore_assignee)
            .ok()
        else {
            continue;
        };

        if !ignore_assignee {
            let Some(jira_assignee) = jira_issue.fields.assignee else {
                debug!("Issue '{}' is unassingned, skip", jira_issue.key);
                continue;
            };

            if let Err(err) = validate_jira_issue(&jira_issue.key, &jira_assignee.email, assignee) {
                debug!("{err}");
                continue;
            }
        }
        debug!(
            "Branch '{local_branch}' with status '{}'",
            jira_issue.fields.status.category.name
        );

        statuses.insert(local_branch, jira_issue.fields.status.category);
    }

    statuses
}

fn extract_issue_name_from_branch(branch_name: &str) -> Result<String> {
    debug!("before extract {branch_name}");

    let issue_name = ISSUE_REGEX
        .find(branch_name)
        .ok_or_else(|| CleanCmdError::ExtractIssueFromBranch(branch_name.to_string()))?
        .as_str()
        .trim();

    let issue_name = match issue_name.strip_suffix(branch::ISSUE_TYPE_SEPARATOR) {
        Some(i) => i,
        None => issue_name,
    };

    Ok(issue_name.to_string())
}

fn validate_jira_issue(issue_key: &str, email: &str, assignee: &str) -> Result<()> {
    let username = extract_username_from_email(email)?;

    if assignee.trim() != username {
        return Err(CleanCmdError::InvalidAssignee {
            issue_key: issue_key.to_string(),
            username,
            assignee: assignee.to_string(),
        });
    }

    Ok(())
}

fn delete_local_branch(branch_name: &str) {
    match git::delete_local_branch(branch_name) {
        Ok(output) => {
            info!("{output}");
        }
        Err(err) => {
            info!("local branch: [{branch_name}] err {err}\n");
        }
    }
}

fn delete_remote_branch(remote: &str, branch_name: &str) {
    match git::delete_remote_branch(remote, branch_name) {
        Ok(output) => {
            info!("{output}");
        }
        Err(err) => {
            info!("remote branch: [{branch_name}] err {err}\n");
        }
    }
}

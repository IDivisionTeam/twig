use anyhow::{Result, bail};
use clap::{Args, Subcommand};
use log::info;
use std::collections::HashSet;

use crate::config::RemoteProvider;
use crate::network::api::HttpClient;
use crate::network::model::JiraIssueType;
use crate::vcs::{RemoteChangesParams, VCSClient};
use crate::{
    branch::{Branch, BranchType},
    config, git,
    network::api::JiraApi,
};

#[derive(Args)]
pub struct Create {
    #[command(subcommand)]
    sub: Subcommands,
}

#[derive(Subcommand)]
pub enum Subcommands {
    // Create a branch for an issue
    Branch {
        issue: String,

        #[arg(short, long)]
        r#type: Option<String>,

        #[arg(short, long, default_value_t = false)]
        push: bool,
    },

    /// Creates changes on remote, for example, a pull request for GitHub provider
    #[clap(visible_aliases = &["pr", "mr"])]
    ChangesOnRemote {
        #[arg(short, long, value_delimiter = ',')]
        labels: Vec<String>,

        #[arg(short, long, default_value_t = false)]
        push: bool,
    },
}

pub fn handle<C: HttpClient>(
    jira_api: &JiraApi<C>,
    vcs_client: &dyn VCSClient,
    args: &Create,
    config: &config::Config,
) -> Result<()> {
    match &args.sub {
        Subcommands::Branch { issue, r#type, push } => {
            create_issue(jira_api, issue, r#type.clone(), *push, config)
        }
        Subcommands::ChangesOnRemote { labels, push } => match &config.remote {
            Some(r) => create_change_on_remote(
                labels.to_owned(),
                &config.project.remote,
                &config.project.branch,
                vcs_client,
                *push,
                r,
            ),
            None => bail!("failed to make changes on remote: remote config is not set"),
        },
    }
}

fn create_issue<C: HttpClient>(
    jira_api: &JiraApi<C>,
    issue: &str,
    issue_type: Option<String>,
    push: bool,
    config: &config::Config,
) -> Result<()> {
    let jira_issue = jira_api.get_jira_issue(issue)?;

    let branch_type = issue_type.or(try_map_issue_type_to_branch_type(
        jira_issue.fields.issue_type.as_ref(),
        &config.mapping,
    ));

    let exclude_phrases = config
        .project
        .exclude_phrases
        .iter()
        .map(String::as_str)
        .collect();
    let b = Branch::new(branch_type, exclude_phrases);

    let branch_name = b.build_name(&jira_issue.key, jira_issue.fields.summary.as_ref());
    let output = git::checkout(&branch_name)?;
    info!("{output}");

    if push {
        let output = git::push_to_remote(&branch_name, &config.project.remote)?;
        info!("{output}");
    }

    Ok(())
}

fn try_map_issue_type_to_branch_type(
    issue_type: Option<&JiraIssueType>,
    mapping: &config::Mapping,
) -> BranchType {
    issue_type
        .as_ref()
        .and_then(|it| mapping.entries.get(&it.id))
        .cloned()
}

fn create_change_on_remote(
    labels: Vec<String>,
    remote: &str,
    default_branch: &str,
    vcs_client: &dyn VCSClient,
    push: bool,
    config: &config::Remote,
) -> Result<()> {
    let branch = git::current_branch()?;
    if push {
        let output = git::push_to_remote(&branch, remote)?;
        info!("{output}");
    } else {
        git::get_remote_branch(&branch)?
            .ok_or_else(|| anyhow::anyhow!("branch `{branch}` has no remote branch"))?;
    }

    let origin_url = git::get_origin_url()?;
    let (owner, repo) = extract_owner_and_repo(&origin_url)?;

    let mut combined_labels = config.labels.clone();
    combined_labels.extend(labels);
    let mut seen = HashSet::new();
    combined_labels.retain(|item| seen.insert(item.clone()));

    let url = vcs_client.create_changes_on_remote(RemoteChangesParams {
        owner: owner.to_string(),
        repo: repo.to_string(),
        branch: branch.clone(),
        default_branch: default_branch.to_string(),
        title: git::latest_commit_msg()?,
        labels: combined_labels,
    })?;

    match config.provider {
        RemoteProvider::GitHub => info!("Pull request created: {url}"),
        RemoteProvider::Gitlab => info!("Merge request created: {url}"),
    }

    Ok(())
}

fn extract_owner_and_repo(origin_url: &str) -> Result<(&str, &str)> {
    let origin_url = origin_url.trim().trim_end_matches(".git").trim_matches('/');

    let origin_url_parts = origin_url.split('/').collect::<Vec<_>>();
    let prefix_with_origin = origin_url_parts[origin_url_parts.len() - 2]
        .split(':')
        .collect::<Vec<_>>();
    let owner = prefix_with_origin
        .last()
        .ok_or_else(|| anyhow::anyhow!("failed to extract owner from origin: {origin_url}"))?;
    let repo = origin_url_parts
        .last()
        .ok_or_else(|| anyhow::anyhow!("failed to extract repo from origin: {origin_url}"))?;

    Ok((owner, repo))
}

#[cfg(test)]
#[allow(clippy::result_large_err)]
mod tests {
    use super::*;
    use crate::git;
    use crate::network::client::ApiError;
    use figment::Jail;
    use rstest::rstest;
    use serde::Serialize;
    use serde::de::DeserializeOwned;
    use std::collections::HashMap;
    use std::sync::{Mutex, Once};

    #[rstest]
    #[case::from_type_arg(Some("test-type".to_string()), "test-type")]
    #[case::from_mapping(None, "test-mapping-type")]
    fn test_handle_branch_type(#[case] branch_type_arg: BranchType, #[case] branch_type_prefix: String) {
        Jail::expect_with(|_| {
            init_logger();
            let issue = "test-issue";
            let issue_type = branch_type_arg;
            let push = false;

            let config = config::Config {
                credentials: config::Credentials::default(),
                project: config::Project::default(),
                remote: None,
                mapping: config::Mapping {
                    entries: HashMap::from([("test-branch-type".to_string(), branch_type_prefix.clone())]),
                },
            };

            init_git_with_commit("Initial commit").unwrap();

            let mut mock_client = MockHttpClient::new();
            let mock_response = serde_json::json!({
                "id": "test-id",
                "key": "test-key",
                "fields": {
                    "issuetype": {
                        "id": "test-branch-type",
                        "name": "issue-type-name"
                    },
                    "summary": "This is a mock summary",
                    "status": null,
                    "assignee": null
                }
            });
            mock_client.set_get_response(mock_response);
            create_issue(&JiraApi::new(mock_client), issue, issue_type, push, &config).unwrap();

            let new_branch = git::execute(&["branch", "--show-current"]).unwrap();

            assert_eq!(
                branch_type_prefix + "/test-key_this-is-mock-summary",
                new_branch.trim()
            );
            Ok(())
        });
    }

    #[rstest]
    fn test_handle_pushed_to_remote() {
        Jail::expect_with(|jail| {
            init_logger();
            let issue = "test-issue";
            let issue_type = None;
            let push = true;

            let config = config::Config {
                credentials: config::Credentials::default(),
                project: config::Project::default(),
                remote: None,
                mapping: config::Mapping {
                    entries: HashMap::new(),
                },
            };

            init_git_with_commit("Initial commit").unwrap();

            create_tmp_remote(jail, "remote")?;

            let mut mock_client = MockHttpClient::new();
            let mock_response = serde_json::json!({
                "id": "test-id",
                "key": "test-key",
                "fields": {
                    "issue_type": null,
                    "summary": "This is a mock summary",
                    "status": null,
                    "assignee": null
                }
            });
            mock_client.set_get_response(mock_response);

            create_issue(&JiraApi::new(mock_client), issue, issue_type, push, &config).unwrap();

            jail.change_dir("remote")?;
            assert!(git::branch_exists("test-key_this-is-mock-summary").unwrap());

            Ok(())
        });
    }

    #[rstest]
    fn test_create_change_on_remote() {
        Jail::expect_with(|jail| {
            init_logger();
            let labels_from_args = vec!["test: label".to_string()];
            let default_branch = "master";
            let current_branch = "test-branch";
            let commit_message = "test commit message";
            let push = true;

            let config = config::Remote::new(
                RemoteProvider::GitHub,
                "test-token".to_string(),
                None,
                vec!["label-from-config: test".to_string()],
            );
            let vcs_client = MockVCSClient {};

            let owner = "owner";
            let repo = "remote";
            let remote_path = format!("{owner}/{repo}");
            init_git_with_commit(commit_message).unwrap();
            git::checkout(current_branch).unwrap();
            create_tmp_remote(jail, &remote_path)?;

            create_change_on_remote(
                labels_from_args,
                "origin",
                default_branch,
                &vcs_client,
                push,
                &config,
            )
            .unwrap();

            jail.change_dir(&remote_path)?;
            assert!(git::branch_exists(current_branch).unwrap());

            let expected_remote_params = RemoteChangesParams {
                owner: owner.to_string(),
                repo: repo.to_string(),
                branch: current_branch.to_string(),
                default_branch: default_branch.to_string(),
                title: commit_message.to_string(),
                labels: vec!["label-from-config: test".to_string(), "test: label".to_string()],
            };

            let logs = LOG_BUFFER.lock().unwrap();
            let actual_params = logs.iter().find(|l| l.contains("Pull request created"));
            assert!(actual_params.is_some());
            assert_eq!(
                format!("{expected_remote_params:?}"),
                actual_params.unwrap().replace("Pull request created: ", "")
            );

            Ok(())
        });
    }

    #[rstest]
    fn test_create_change_on_remote_error_if_no_push_and_no_remote_branch() {
        Jail::expect_with(|_| {
            init_logger();
            let labels_from_args = vec!["test: label".to_string()];
            let default_branch = "master";
            let current_branch = "test-branch";
            let commit_message = "test commit message";
            let push = false;

            let config = config::Remote::new(
                RemoteProvider::GitHub,
                "test-token".to_string(),
                None,
                vec!["label-from-config: test".to_string()],
            );

            init_git_with_commit(commit_message).unwrap();
            git::checkout(current_branch).unwrap();

            let result = create_change_on_remote(
                labels_from_args,
                "origin",
                default_branch,
                &MockVCSClient {},
                push,
                &config,
            );
            assert_eq!(
                "branch `test-branch` has no remote branch",
                result.unwrap_err().to_string()
            );

            Ok(())
        });
    }

    fn init_git_with_commit(commit_msg: &str) -> Result<()> {
        git::execute(&["init"])?;
        git::execute(&[
            "-c",
            "user.name=test@example.com",
            "-c",
            "user.email=test name",
            "commit",
            "--allow-empty",
            "-m",
            commit_msg,
        ])?;
        Ok(())
    }

    fn create_tmp_remote(jail: &mut Jail, remote_path: &str) -> std::result::Result<(), figment::Error> {
        let remote_dir = jail.create_dir(remote_path)?;
        git::execute(&[
            "remote",
            "add",
            "origin",
            &format!(
                "{}/.git",
                jail.directory()
                    .join(&remote_dir)
                    .into_os_string()
                    .into_string()
                    .unwrap()
            ),
        ])
        .unwrap();

        jail.change_dir(&remote_dir)?;
        git::execute(&["init"]).unwrap();
        jail.change_dir(jail.directory())?;

        Ok(())
    }

    struct MockHttpClient {
        get_response_json: serde_json::Value,
    }

    impl MockHttpClient {
        fn new() -> MockHttpClient {
            MockHttpClient {
                get_response_json: serde_json::Value::Null,
            }
        }

        fn set_get_response(&mut self, response_json: serde_json::Value) {
            self.get_response_json = response_json;
        }

        fn request<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
            result: &serde_json::Value,
        ) -> Result<T> {
            let result: T = serde_json::from_value(result.clone())
                .map_err(|e| anyhow::anyhow!("Mock deserialization failed: {e}"))?;
            Ok(result)
        }
    }

    impl HttpClient for MockHttpClient {
        fn get<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
            &self,
            _: &str,
            _: Vec<(&str, &str)>,
        ) -> Result<T> {
            Self::request::<T, E>(&self.get_response_json)
        }

        fn post<T: DeserializeOwned, S: Serialize, E: ApiError + DeserializeOwned>(
            &self,
            _: &str,
            _: S,
        ) -> Result<T> {
            todo!()
        }
    }

    struct MockVCSClient {}

    impl VCSClient for MockVCSClient {
        fn create_changes_on_remote(&self, params: RemoteChangesParams) -> Result<String> {
            Ok(format!("{params:?}"))
        }
    }

    static INITIALIZER: Once = Once::new();
    static LOG_BUFFER: Mutex<Vec<String>> = Mutex::new(Vec::new());

    fn init_logger() {
        INITIALIZER.call_once(|| {
            log::set_logger(&SpyLogger).unwrap();
            log::set_max_level(log::LevelFilter::Info);
        });
        LOG_BUFFER.lock().unwrap().clear();
    }

    struct SpyLogger;
    impl log::Log for SpyLogger {
        fn enabled(&self, _: &log::Metadata) -> bool {
            true
        }
        fn log(&self, record: &log::Record) {
            let msg = format!("{}", record.args());
            eprintln!("[{}] {}", record.level(), msg);
            if let Ok(mut buf) = LOG_BUFFER.lock() {
                buf.push(msg);
            }
        }
        fn flush(&self) {}
    }
}

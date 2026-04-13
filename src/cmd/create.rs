use anyhow::Result;
use clap::Args;
use log::info;

use crate::network::api::HttpClient;
use crate::network::model::JiraIssueType;
use crate::{branch::{BranchType, Branch}, config, git, network::api::JiraApi};

#[derive(Args)]
pub struct Create {
    issue: String,

    #[arg(short, long)]
    r#type: Option<String>,

    #[arg(short, long, default_value_t = false)]
    push: bool,
}

pub fn handle<C: HttpClient>(
    jira_api: &JiraApi<C>,
    args: &Create,
    config: &config::Config,
) -> Result<()> {
    let jira_issue = jira_api.get_jira_issue(&args.issue)?;

    let branch_type: BranchType = args.r#type.clone().or(try_map_issue_type_to_branch_type(
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

    if args.push {
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

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use super::*;
    use crate::git;
    use crate::network::client::ApiError;
    use figment::Jail;
    use rstest_log::rstest;
    use serde::Serialize;
    use serde::de::DeserializeOwned;

    #[rstest]
    #[case::from_type_arg(Some("test-type".to_string()), "test-type")]
    #[case::from_mapping(None, "test-mapping-type")]
    fn test_handle_branch_type(
        #[case] branch_type_arg: BranchType,
        #[case] branch_type_prefix: String,
    ) {
        Jail::expect_with(|_| {
            let args = Create {
                issue: "test-issue".to_string(),
                r#type: branch_type_arg,
                push: false,
            };

            let config = config::Config {
                credentials: config::Credentials::default(),
                project: config::Project::default(),
                mapping: config::Mapping {
                    entries: HashMap::from([(
                        "test-branch-type".to_string(),
                        branch_type_prefix.to_string(),
                    )]),
                },
            };

            git::execute(&["config", "--global", "user.email", "test@example.com"]).unwrap();
            git::execute(&["config", "--global", "user.name", "test name"]).unwrap();

            git::execute(&["init"]).unwrap();
            git::execute(&["commit", "--allow-empty", "-m", "Initial commit"]).unwrap();

            let mut mock_client = MockClient::new();
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
            handle(&JiraApi::new(mock_client), &args, &config).unwrap();

            let new_branch = git::execute(&["branch", "--show-current"]).unwrap();

            assert_eq!(
                branch_type_prefix + "/test-key_this-is-mock-summary",
                new_branch.trim()
            );
            Ok(())
        })
    }

    #[rstest]
    fn test_handle_pushed_to_remote() {
        Jail::expect_with(|jail| {
            let args = Create {
                issue: "test-issue".to_string(),
                r#type: None,
                push: true,
            };

            let config = config::Config {
                credentials: config::Credentials::default(),
                project: config::Project::default(),
                mapping: config::Mapping {
                    entries: HashMap::new(),
                },
            };

            git::execute(&["config", "--global", "user.email", "test@example.com"]).unwrap();
            git::execute(&["config", "--global", "user.name", "test name"]).unwrap();

            git::execute(&["init"]).unwrap();
            git::execute(&["commit", "--allow-empty", "-m", "Initial commit"]).unwrap();

            let remote_dir = jail.create_dir("remote")?;
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

            let mut mock_client = MockClient::new();
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

            handle(&JiraApi::new(mock_client), &args, &config).unwrap();

            jail.change_dir(&remote_dir)?;
            assert_eq!(
                true,
                git::branch_exists("test-key_this-is-mock-summary").unwrap()
            );

            Ok(())
        })
    }

    struct MockClient {
        get_response_json: serde_json::Value,
    }

    impl MockClient {
        fn new() -> MockClient {
            MockClient {
                get_response_json: serde_json::Value::Null,
            }
        }

        fn set_get_response(&mut self, response_json: serde_json::Value) {
            self.get_response_json = response_json;
        }

        fn request<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
            &self,
            result: &serde_json::Value,
        ) -> Result<T> {
            let result: T = serde_json::from_value(result.clone())
                .map_err(|e| anyhow::anyhow!("Mock deserialization failed: {}", e))?;
            Ok(result)
        }
    }

    impl HttpClient for MockClient {
        fn get<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
            &self,
            _: &str,
            _: Vec<(&str, &str)>,
        ) -> Result<T> {
            self.request::<T, E>(&self.get_response_json)
        }

        fn post<T: DeserializeOwned, S: Serialize, E: ApiError + DeserializeOwned>(
            &self,
            _: &str,
            _: S,
        ) -> Result<T> {
            todo!()
        }
    }
}

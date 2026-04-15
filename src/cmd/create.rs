use anyhow::{bail, Context, Result};
use clap::{Args, Subcommand};
use log::info;

use crate::network::api::HttpClient;
use crate::network::model::JiraIssueType;
use crate::{branch::{BranchType, Branch}, config, git, network::api::JiraApi};
use crate::vcs::VCSClient;

#[derive(Args)]
pub struct Create {
    issue: Option<String>,

    #[arg(short, long)]
    r#type: Option<String>,

    #[arg(short, long, default_value_t = false)]
    push: bool,

    #[command(subcommand)]
    sub: Option<Subcommands>
}


#[derive(Subcommand)]
pub enum Subcommands {

    #[clap(visible_aliases = &["pr", "mr"])]
    ChangesOnRemote {
        #[arg(short, long, value_parser = parse_key_val)]
        labels: Vec<(String, String)>,
    }
}

fn parse_key_val(s: &str) -> Result<(String, String)> {
    s.split_once('=')
        .map(|(k, v)| (k.to_string(), v.to_string()))
        .context(format!("invalid key=value: no '=' found in '{s}'"))
}

pub fn handle<C: HttpClient>(
    jira_api: &JiraApi<C>,
    vcs_client: Box<dyn VCSClient>,
    args: &Create,
    config: &config::Config,
) -> Result<()> {
    if let Some(sub) = &args.sub {
        if args.issue.is_some() {
            bail!("only one positional argument can be provided: 'issue' or subcommand ")
        }
        match sub {
            Subcommands::ChangesOnRemote { labels } => {
                create_change_on_remote(labels.as_ref(), vcs_client, config)?;
            }
        }
    } else if let Some(issue) = &args.issue {
        create_issue(jira_api, issue, args.r#type.clone(), args.push, config)?;
    }  else {
        bail!("either positional argument 'issue' or subcommand must be provided")
    }

    Ok(())
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

fn create_change_on_remote(_labels: &[(String, String)],  vcs_client: Box<dyn VCSClient>, _config: &config::Config) -> Result<()> {
    vcs_client.create_changes_on_remote()?;
    // create pr using client
    todo!()
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
            let issue = "test-issue";
            let issue_type = branch_type_arg;
            let push = false;

            let config = config::Config {
                credentials: config::Credentials::default(),
                project: config::Project::default(),
                remote: None,
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
            create_issue(&JiraApi::new(mock_client), issue, issue_type, push, &config).unwrap();

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

            create_issue(&JiraApi::new(mock_client), issue, issue_type, push, &config).unwrap();

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

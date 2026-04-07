use anyhow::Result;
use clap::Args;
use log::info;

use crate::{branch::*, config, git, network::api::JiraApi};
use crate::network::api::HttpClient;
use crate::network::model::JiraIssueType;

#[derive(Args)]
pub struct Create {
    issue: String,

    #[arg(short, long)]
    _type: Option<String>,

    #[arg(short, long, default_value_t = false)]
    push: bool,
}

pub fn handle<C: HttpClient>(jira_api: &JiraApi<C>, args: &Create, config: &config::Config) -> Result<()> {
    let jira_issue = jira_api.get_jira_issue(&args.issue)?;

    let branch_type: BranchType = args._type.clone().or(try_map_issue_type_to_branch_type(&jira_issue.fields.issue_type, &config.mapping));

    let exclude_phrases = config
        .project
        .exclude_phrases
        .iter()
        .map(String::as_str)
        .collect();
    let b = Branch::new(branch_type, exclude_phrases);

    let branch_name = b.build_name(&jira_issue.key, &jira_issue.fields.summary);
    let output = git::checkout(&branch_name)?;
    info!("{output}");

    if args.push {
        let output = git::push_to_remote(&branch_name, &config.project.remote)?;
        info!("{output}");
    }

    Ok(())
}

fn try_map_issue_type_to_branch_type(
    issue_type: &Option<JiraIssueType>,
    mapping: &config::Mapping,
) -> BranchType {
    issue_type.as_ref().and_then(|it|  mapping.entries.get(&it.id)).cloned()
}

//
// #[cfg(test)]
// mod tests {
//     use serde::de::DeserializeOwned;
//     use serde::Serialize;
//     use crate::network::client::ApiError;
//     use super::*;
//
//     struct MockClient;
//
//     impl HttpClient for MockClient {
//         fn get<T: DeserializeOwned, E: ApiError + DeserializeOwned>(&self, path: &str, params: Vec<(&str, &str)>) -> Result<T> {
//             todo!()
//         }
//
//         fn post<T: DeserializeOwned, S: Serialize, E: ApiError + DeserializeOwned>(&self, path: &str, body: S) -> Result<T> {
//             todo!()
//         }
//     }
//
//     #[test]
//     fn test_handle() {
//
//     }
//
// }
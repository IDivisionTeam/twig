use anyhow::Result;
use clap::Args;
use log::info;

use crate::{
    branch::*,
    config, git,
    network::{api::JiraApi, model::JiraIssueType},
};

#[derive(Args)]
pub struct Create {
    issue: String,

    #[arg(short, long)]
    _type: Option<String>,

    #[arg(short, long, default_value_t = false)]
    push: bool,
}

pub fn handle(jira_api: &JiraApi, args: &Create, config: &config::Config) -> Result<()> {
    let jira_issue = jira_api.get_jira_issue(&args.issue)?;

    let branch_type = match &args._type {
        Some(_type) => _type.parse()?,
        None => try_map_issue_type_to_branch_type(&jira_issue.fields.issue_type, &config.mapping)?,
    };

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
) -> Result<BranchType> {
    let issue_type = match issue_type {
        Some(it) => it,
        None => return Ok(BranchType::Unspecified),
    };

    match mapping.get(&issue_type.id) {
        Some(mt) => Ok((*mt).into()),
        None => Ok(BranchType::Unspecified),
    }
}

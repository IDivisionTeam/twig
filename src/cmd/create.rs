use anyhow::Result;
use clap::Args;
use log::info;

use crate::{branch::*, config, git, network::api::JiraApi};

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

    let branch_type: BranchType = args._type.clone();

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

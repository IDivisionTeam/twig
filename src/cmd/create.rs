use anyhow::Result;
use clap::Args;

use crate::network::api::JiraApi;

#[derive(Args)]
pub struct Create {
    issue: String,

    #[arg(short, long)]
    _type: Option<String>,

    #[arg(short, long, default_value_t = false)]
    push: bool,
}

pub fn handle(jira_api: &JiraApi, args: &Create) -> Result<()> {
    dbg!(jira_api.get_jira_issue_status(args.issue.clone(), false)?);
    dbg!(jira_api.get_jira_issue_statuses(vec!["10001".to_string(), "10008".to_string()], false)?);
    Ok(())
}

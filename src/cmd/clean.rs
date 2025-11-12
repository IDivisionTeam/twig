use anyhow::Result;
use clap::{Arg, Args, Command};

use crate::network::api::JiraApi;

#[derive(Args)]
pub struct Clean {
    #[arg(short, long)]
    assignee: Option<String>,

    #[arg(short, long, default_value_t = false)]
    any: bool,
}

// FIXME: assignee must default to project.email.
pub fn handle(jira_api: &JiraApi, args: &Clean) -> Result<()> {
    todo!();
    Ok(())
}

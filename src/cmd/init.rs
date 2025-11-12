use anyhow::Result;
use clap::Args;

use crate::network::api::JiraApi;

#[derive(Args)]
pub struct Init;

pub fn handle(jira_api: &JiraApi, args: &Init) -> Result<()> {
    todo!();
    Ok(())
}

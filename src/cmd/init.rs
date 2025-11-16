use anyhow::Result;
use clap::Args;

use crate::{config, network::api::JiraApi};

#[derive(Args)]
pub struct Init {
    /// create a global config
    #[arg(long, default_value_t = false)]
    global: bool,
}

pub fn handle(jira_api: &JiraApi, args: &Init) -> Result<()> {
    config::create_config_if_not_exist(args.global)?;
    Ok(())
}

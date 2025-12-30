use anyhow::Result;
use clap::{Args, Subcommand};

use crate::network::api::JiraApi;

#[derive(Args)]
pub struct Cfg {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    List,
    Get { name: String },
    Set { name: String, value: String },
}

pub fn handle(_jira_api: &JiraApi, _args: &Cfg) -> Result<()> {
    todo!()
}

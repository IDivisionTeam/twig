use anyhow::Result;
use clap::Args;

use crate::{config, network::api::JiraApi};
use crate::network::api::HttpClient;

#[derive(Args)]
pub struct Init {
    /// create a global config
    #[arg(long, default_value_t = false)]
    global: bool,
}

pub fn handle<C: HttpClient>(_: &JiraApi<C>, args: &Init) -> Result<()> {
    let config_path = if args.global {
        config::get_config_global_path()
    } else {
        config::get_config_local_path()
    };
    config::create_config_if_not_exist(&config_path)?;
    Ok(())
}

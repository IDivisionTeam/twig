mod branch;
mod cmd;
mod config;
mod git;
mod network;

use anyhow::Result;
use cmd::twig;

use crate::config::Config;

fn main() -> Result<()> {
    env_logger::Builder::from_default_env()
        .filter_level(log::LevelFilter::Info)
        .init();

    let config: Config = config::read_config()?;
    twig::execute(config)
}

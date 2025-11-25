mod branch;
mod cmd;
mod config;
mod network;

use anyhow::Result;
use cmd::twig;

use crate::config::Config;

fn main() -> Result<()> {
    let config: Config = config::read_config()?;
    twig::execute(config)
}

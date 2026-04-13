mod branch;
mod cmd;
mod config;
mod git;
mod network;


use anyhow::Result;
use cmd::twig;
use env_logger::Env;

use crate::config::Config;

fn main() -> Result<()> {
    let env = Env::new().filter_or("RUST_LOG", "info");
    env_logger::init_from_env(env);

    let config: Config = config::read_config()?;
    twig::execute(&config)
}

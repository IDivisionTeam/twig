use anyhow::{Context, Result};
use clap::{Parser, Subcommand};

use crate::cmd::{cfg, clean, create, init};
use crate::config::Config;
use crate::network;
use crate::network::api::JiraApi;

// To add styling, see https://docs.rs/clap/latest/clap/_derive/_cookbook/cargo_example_derive/index.html
#[derive(Parser)]
#[command(version, about, long_about = None, arg_required_else_help(true))]
struct Twig {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    // FIXME: Add missing descriptions (help).
    /// Deletes branches which have Jira tickets in 'Done' state
    Clean(clean::Clean),
    /// You can query/set/replace options with this command.
    /// The name is the section and the key separated by a dot.
    Config(cfg::Cfg),
    /// create an issue
    Create(create::Create),
    /// blablabla
    Init(init::Init),
}

pub fn execute(config: Config) -> Result<()> {
    let credentials = network::model::Credentials {
        host: config.credentials.host.clone(),
        email: config.credentials.email.clone(),
        auth: config.credentials.auth.clone(),
        token: config.credentials.token.clone(),
    };
    let jira_api = JiraApi::new(credentials)?;

    let twig = Twig::parse();

    match &twig.command {
        Commands::Clean(args) => clean::handle(&jira_api, args, &config).context("clean command failed")?,
        Commands::Config(args) => cfg::handle(args, &config).context("config command failed")?,
        Commands::Create(args) => create::handle(&jira_api, args, &config).context("create command failed")?,
        Commands::Init(args) => init::handle(&jira_api, args).context("init command failed")?,
    }

    Ok(())
}

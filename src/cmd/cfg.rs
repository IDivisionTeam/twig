use crate::config;
use crate::config::read_config_value;
use clap::{Args, Subcommand};
use figment::value::Value;
use thiserror::Error;

#[derive(Args)]
pub struct Cfg {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// List all variables set in config file, along with their values.
    List,
    /// Emits the value of the specified key.
    Get {
        name: String,
    },
    Set {
        name: String,
        value: String,
    },
}

type Result<T> = std::result::Result<T, ConfigCmdError>;

#[derive(Error, Debug, PartialEq, Eq)]
pub enum ConfigCmdError {
    #[error("failed to get value by key {0}")]
    InvalidKey(String),
    #[error("failed to deserialize value")]
    DeserializeValue(),
}

pub fn handle(args: &Cfg, config: &config::Config) -> Result<()> {
    match &args.command {
        Commands::List => handle_list_cmd(config),
        Commands::Get { name } => handle_get_cmd(&name)?,
        #[allow(unused_variables)]
        Commands::Set { name, value } => {
            todo!()
        }
    }

    Ok(())
}

fn handle_list_cmd(config: &config::Config) {
    print!(
        "{}\n{}\n{}",
        config.credentials, config.project, config.mapping
    );
}

fn handle_get_cmd(name: &str) -> Result<()> {
    let value: Value = read_config_value(name)
        .unwrap()
        .ok_or_else(|| ConfigCmdError::InvalidKey(name.to_string()))?;

    let value_to_print: Option<String> = value.deserialize().ok();
    match value_to_print {
        Some(v) => {
            println!("{}", v);
            Ok(())
        }
        None => Err(ConfigCmdError::DeserializeValue()),
    }
}

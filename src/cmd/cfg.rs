use crate::config;
use crate::config::read_config_value;
use clap::{Args, Subcommand};
use figment::value::{Empty, Value};
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

#[derive(Error, Debug)]
pub enum ConfigCmdError {
    #[error("failed to get value by key {0}")]
    InvalidKey(String),
    #[error("failed to get any results")]
    UnexpectedError,
    #[error("failed to read config")]
    ConfigError(#[from] config::ConfigError),
}

pub fn handle(args: &Cfg, config: &config::Config) -> Result<()> {
    match &args.command {
        Commands::List => handle_list_cmd(config),
        Commands::Get { name } => handle_get_cmd(name)?,
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
    let value: Value =
        read_config_value(name)?.ok_or(ConfigCmdError::InvalidKey(name.to_string()))?;

    let value_to_print = print_value(&value);
    match value_to_print {
        Some(v) => {
            println!("{}", v);
            Ok(())
        }

        None => Err(ConfigCmdError::UnexpectedError),
    }
}

fn print_value(value: &Value) -> Option<String> {
    match value {
        Value::Empty(_, e) => match e {
            Empty::None => Some("unspecified".to_string()),
            Empty::Unit => None,
        },
        Value::Bool(_, b) => Some(b.to_string()),
        Value::Num(_, n) => Some(n.to_actual().to_string()),
        Value::String(_, s) => Some(s.clone()),
        Value::Char(_, c) => Some(c.to_string()),
        Value::Array(_, arr) => Some(
            arr.iter()
                .filter_map(print_value)
                .collect::<Vec<_>>()
                .join(", "),
        ),
        Value::Dict(_, map) => Some(
            map.iter()
                .filter_map(|(k, v)| {
                    let value = print_value(v).unwrap_or(String::new());
                    Some(format!("{} = {}", k, value))
                })
                .collect::<Vec<_>>()
                .join("\n"),
        ),
    }
}

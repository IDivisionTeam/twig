use crate::config;
use anyhow::Result;
use clap::{Args, Subcommand};

#[derive(Args)]
pub struct Cfg {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// List all variables set in config file, along with their values.
    List,
    Get {
        name: String,
    },
    Set {
        name: String,
        value: String,
    },
}

pub fn handle(args: &Cfg, config: &config::Config) -> Result<()> {
    match &args.command {
        Commands::List => handle_list_cmd(config),
        #[allow(unused_variables)]
        Commands::Get { name } => {
            todo!()
        }
        #[allow(unused_variables)]
        Commands::Set { name, value } => {
            todo!()
        }
    }

    Ok(())
}

fn handle_list_cmd(config: &config::Config) {
    print!("{}\n{}\n{}", config.credentials, config.project, config.mapping);
}

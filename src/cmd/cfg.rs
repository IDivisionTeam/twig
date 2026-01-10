use crate::config;
use crate::config::MappingType;
use anyhow::Result;
use clap::{Args, Subcommand};
use std::collections::BTreeMap;

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
    print!("{}{}", config.credentials, config.project);

    let inverted_mapping = config.mapping.clone().into_iter().fold(
        BTreeMap::<MappingType, Vec<String>>::new(),
        |mut acc, (k, v)| {
            acc.entry(v).or_default().push(k);
            acc
        },
    );

    let size = inverted_mapping.len().saturating_sub(1);
    for (index, (key, value)) in inverted_mapping.into_iter().enumerate() {
        let values = value.join(", ");
        let output = format!("mapping.{key}=[{values}]");

        if index >= size {
            print!("{output}");
        } else {
            println!("{output}");
        }
    }
}

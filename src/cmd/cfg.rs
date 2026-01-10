use crate::config;
use crate::config::MappingType;
use anyhow::Result;
use clap::{Args, Subcommand};
use std::collections::{BTreeMap, HashMap};

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
        Commands::List => handle_list_cmd(&config),
        Commands::Get { name } => {
            todo!()
        }
        Commands::Set { name, value } => {
            todo!()
        }
    }

    Ok(())
}

fn handle_list_cmd(config: &config::Config) {
    let output = format!(
        "\
        credentials.auth={auth}\n\
        credentials.email={email}\n\
        credentials.host={host}\n\
        credentials.token={token}\n\
        project.branch={branch}\n\
        project.exclude_phrases=[{exclude_phrases}]\n\
        project.remote={remote}\n\
        ",
        auth = config.credentials.auth,
        host = config.credentials.host,
        email = config.credentials.email,
        token = config.credentials.token,
        branch = config.project.branch,
        exclude_phrases = config.project.exclude_phrases.join(", "),
        remote = config.project.remote,
    );

    print!("{output}");

    let inverted_mapping = config.mapping.clone().into_iter().fold(
        HashMap::<MappingType, Vec<String>>::new(),
        |mut acc, (k, v)| {
            acc.entry(v).or_default().push(k);
            acc
        },
    );

    let sorted_mapping: BTreeMap<_, _> = inverted_mapping.into_iter().collect();

    let mut index = 0;
    let size = sorted_mapping.len() - 1;
    for (key, value) in sorted_mapping {
        let values = value.join(", ");
        let output = format!("mapping.{key}=[{values}]");

        if index >= size {
            print!("{output}");
        } else {
            println!("{output}");
        }

        index += 1;
    }
}

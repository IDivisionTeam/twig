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
    Get { name: String },
    Set { name: String, value: String },
}

pub fn handle(args: &Cfg) -> Result<()> {
  todo!()
}

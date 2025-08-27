use clap::Command;

pub const NAME: &str = "init";

pub fn subcommand() -> Command {
    Command::new(NAME)
}
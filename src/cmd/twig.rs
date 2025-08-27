use clap::builder::Styles;
use clap::{command, Command};

use crate::{cmd::cfg, cmd::clean, cmd::create, cmd::init};

fn twig() -> Command {
    command!()
        .styles(Styles::styled())
        .arg_required_else_help(true)
        // FIXME: Add missing descriptions (help).
        .subcommand(clean::subcommand())
        .subcommand(cfg::subcommand())
        .subcommand(create::subcommand())
        .subcommand(init::subcommand())
}

pub fn execute() {
    let matches = twig().get_matches();

    // TODO: handle matches
    match matches.subcommand() {
        Some((clean::NAME, _clean_matches)) => {
            // FIXME: assignee must default to project.email.
        }
        Some((cfg::NAME, _config_matches)) => {}
        Some((create::NAME, _create_matches)) => {}
        Some((init::NAME, _init_matches)) => {}
        _ => unreachable!(), // all commands are defined above, anything else is unreachable!()
    }
}

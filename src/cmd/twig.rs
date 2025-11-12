use anyhow::Result;
use clap::builder::Styles;
use clap::{Command, command};

use crate::network;
use crate::network::api::JiraApi;
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

pub fn execute() -> Result<()> {
    let credentials = network::model::Credentials {
        host: "".to_string(),
        auth: "".to_string(),
        email: "".to_string(),
        token: "".to_string(),
    };
    let jira_api = JiraApi::new(credentials)?;

    let matches = twig().get_matches();

    // TODO: handle matches
    match matches.subcommand() {
        Some((clean::NAME, _clean_matches)) => {
            // FIXME: assignee must default to project.email.
        }
        Some((cfg::NAME, _config_matches)) => {}
        Some((create::NAME, _create_matches)) => {
            let issue_key = _create_matches
                .get_one::<String>("issue")
                .ok_or(anyhow::Error::msg("issues expected in arg"))?;
            dbg!(jira_api.get_jira_issue_status(issue_key.clone(), false)?);
            dbg!(
                jira_api.get_jira_issue_statuses(
                    vec!["10001".to_string(), "10008".to_string()],
                    false
                )?
            );
        }
        Some((init::NAME, _init_matches)) => {}
        _ => unreachable!(), // all commands are defined above, anything else is unreachable!()
    }
    Ok(())
}

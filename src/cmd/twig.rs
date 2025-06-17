use clap::builder::{NonEmptyStringValueParser, Styles};
use clap::{command, value_parser, Arg, Command};

const CLEAN: &str = "clean";
const CONFIG: &str = "config";
const CREATE: &str = "create";
const INIT: &str = "init";

fn twig() -> Command {
    command!()
        .styles(Styles::styled())
        .arg_required_else_help(true)
        // FIXME: Add missing descriptions (help).
        .subcommand(clean_subcommand())
        .subcommand(config_subcommand())
        .subcommand(create_subcommand())
        .subcommand(init_subcommand())
}

fn clean_subcommand() -> Command {
    Command::new(CLEAN)
        .subcommand(Command::new("all").args(clean_common_args()))
        .subcommand(Command::new("local").args(clean_common_args()))
}

fn clean_common_args() -> Vec<Arg> {
    vec![
        Arg::new("assignee")
            .long("assignee")
            .short('a')
            .num_args(1)
            .value_parser(value_parser!(String))
            .default_value(""),
        Arg::new("any")
            .long("any")
            .num_args(0)
            .value_parser(value_parser!(bool))
            .default_missing_value("false")
            .default_value("false"),
    ]
}

fn config_subcommand() -> Command {
    Command::new(CONFIG)
        .subcommand(Command::new("list"))
        .subcommand(
            Command::new("get").arg(
                Arg::new("name")
                    .num_args(1)
                    .required(true)
                    .value_parser(NonEmptyStringValueParser::new()),
            ),
        )
        .subcommand(
            Command::new("set")
                .arg(
                    Arg::new("name")
                        .num_args(1)
                        .required(true)
                        .value_parser(NonEmptyStringValueParser::new()),
                )
                .arg(
                    Arg::new("value")
                        .num_args(1)
                        .required(true)
                        .value_parser(NonEmptyStringValueParser::new()),
                ),
        )
}

fn create_subcommand() -> Command {
    Command::new(CREATE)
        .arg(
            Arg::new("issue")
                .num_args(1)
                .required(true)
                .value_parser(NonEmptyStringValueParser::new()),
        )
        .arg(
            Arg::new("type")
                .long("type")
                .short('t')
                .num_args(1)
                .value_parser(value_parser!(String))
                .default_value(""),
        )
        .arg(
            Arg::new("push")
                .long("push")
                .short('p')
                .num_args(0)
                .value_parser(value_parser!(bool))
                .default_missing_value("false")
                .default_value("false"),
        )
}

fn init_subcommand() -> Command {
    Command::new(INIT)
}

pub fn execute() {
    let matches = twig().get_matches();

    // TODO: handle matches
    match matches.subcommand() {
        Some((CLEAN, _clean_matches)) => {
            // FIXME: assignee must default to project.email.
        }
        Some((CONFIG, _config_matches)) => {}
        Some((CREATE, _create_matches)) => {}
        Some((INIT, _init_matches)) => {}
        _ => unreachable!(), // all commands are defined above, anything else is unreachable!()
    }
}

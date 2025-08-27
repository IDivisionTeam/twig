use clap::{value_parser, Arg, Command};

pub const NAME: &str = "clean";

pub fn subcommand() -> Command {
    Command::new(NAME)
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

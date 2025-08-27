use clap::{value_parser, Arg, Command};
use clap::builder::NonEmptyStringValueParser;

pub const NAME: &str = "create";

pub fn subcommand() -> Command {
    Command::new(NAME)
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
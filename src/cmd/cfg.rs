use clap::{Arg, Command};
use clap::builder::NonEmptyStringValueParser;

pub const NAME: &str = "config";

pub fn subcommand() -> Command {
    Command::new(NAME)
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

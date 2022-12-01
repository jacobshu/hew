use clap::{arg, command, Command};
use log::{error, info, warn};
use pretty_env_logger;
mod load;
mod utils;

fn main() {
    pretty_env_logger::init();

    let matches = command!() // requires `cargo` feature
        .propagate_version(true)
        .subcommand_required(true)
        .arg_required_else_help(true)
        .subcommand(
            Command::new("load")
                .about("Initializes and updates the system")
                .arg(arg!(update: -u --update "run updates only"))
                .arg(arg!(link: -s --link "symlink dotfiles")),
        )
        .get_matches();

    match matches.subcommand() {
        Some(("load", sub_matches)) => {
            load::init(sub_matches);
        }
        _ => error!("Exhausted list of subcommands and subcommand_required prevents `None`"),
    }
}

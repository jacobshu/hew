use clap::{arg, command, Command};
use log::error;
use pretty_env_logger;
use std::env::set_var;
mod load;
mod task;
mod utils;

fn main() {
    set_var("RUST_LOG", "INFO");
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
        .subcommand(
            Command::new("task")
                .about("Start the task manager")
                .arg(arg!(add: -a --add "add a task to the list"))
                .arg(arg!(list: -l --list "list all open tasks")),
        )
        .get_matches();

    match matches.subcommand() {
        Some(("load", sub_matches)) => {
            load::init(sub_matches);
        }
        Some(("task", sub_matches)) => {
            println!("run task with matches: {:?}", sub_matches);
            let _t = task::start();
        }
        _ => error!("Exhausted list of subcommands and subcommand_required prevents `None`"),
    }
}

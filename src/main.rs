use clap::{arg, command, Command};
mod load;

fn main() {
    let matches = command!() // requires `cargo` feature
    .propagate_version(true)
    .subcommand_required(true)
    .arg_required_else_help(true)
    .subcommand(
            Command::new("load")
            .about("Initializes and updates the system")
            .arg(arg!(update: -u --update "run updates only")),
    )
    .get_matches();

    match matches.subcommand() {
        Some(("load", sub_matches)) => println!(
                "load called with arg: {:?}",
            sub_matches.get_flag("update")
        ),
        _ => unreachable!("Exhausted list of subcommands and subcommand_required prevents `None`"),
    }
    load::init();
}
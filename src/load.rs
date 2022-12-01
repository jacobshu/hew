use clap::ArgMatches;
use log::info;
use std::fs;
use std::io;
use std::process::{Command, Stdio};
use crate::utils::which;


pub fn init(args: &ArgMatches) {
    info!(
        "beginning load::  update: {:?} | link: {:?}",
        args.get_flag("update"),
        args.get_flag("link")
    );
    install_homebrew();
}

fn install_homebrew() {
//    if which("brew") {
//        info!("homebrew already installed")
//    } else {
        info!("installing homebrew...");

        let resp = reqwest::blocking::get(
            "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh",
        )
        .expect("request failed");
        let body = resp.text().expect("body invalid");
        let mut out = fs::File::create("homebrew_install.sh").expect("failed to create file");
        io::copy(&mut body.as_bytes(), &mut out).expect("failed to copy content");

        info!("homebrew install script downloaded");
//        let install_script = fs::read_to_string("homebrew_install.sh")
//            .expect("Should have been able to read the file");

        Command::new("/bin/bash")
            .env("NONINTERACTIVE", "1")
            .arg("-c")
            .arg("homebrew_install.sh")
            .spawn()
            .expect("homebrew installation failed to start");

        fs::remove_file("homebrew_install.sh");
//    }
}

//fn link_dotfiles() {}


use clap::ArgMatches;
use log::info;
use std::fs;
use std::io;
use std::process::{Command, Stdio};

pub fn init(args: &ArgMatches) {
    info!(
        "beginning load::  update: {:?} | link: {:?}",
        args.get_flag("update"),
        args.get_flag("link")
    );
    install_homebrew();
}

fn install_homebrew() {
    let mut homebrew_status = Command::new("which")
        .arg("dot")
        .stdout(Stdio::piped())
        .spawn()
        .expect("failed to check homebrew");
    let stdout = homebrew_status.stdout.take().unwrap();
    println!("out: {:?}", stdout);

    let which_output = homebrew_status
        .wait_with_output()
        .expect("failed in wait for 'which'");
    if which_output.status.success() {
        info!("homebrew already installed")
    } else {
        info!("installing homebrew");

        let resp = reqwest::blocking::get(
            "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh",
        )
        .expect("request failed");
        let body = resp.text().expect("body invalid");
        let mut out = fs::File::create("homebrew_install.sh").expect("failed to create file");
        io::copy(&mut body.as_bytes(), &mut out).expect("failed to copy content");
        let install_script = fs::read_to_string("homebrew_install.sh")
            .expect("Should have been able to read the file");

        Command::new("/bin/bash")
        .env("NONINTERACTIVE", "1")
        .arg("-c")
        .arg(install_script)
        .spawn()
        .expect("homebrew installation failed to start");
    }
}

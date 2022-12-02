use clap::ArgMatches;
use log::{info, error};
use std::fs::{ File, read_to_string, remove_file };
use std::io::copy;
use std::process::Command;
use std::path::Path;
use std::os::unix::fs::symlink;
use crate::utils::which;



pub fn init(args: &ArgMatches) {
    info!(
        "beginning load::  update: {:?} | link: {:?}",
        args.get_flag("update"),
        args.get_flag("link")
    );
    install_homebrew();
    link_dotfiles();
}

fn install_homebrew() {
    if which("brew") {
        info!("homebrew already installed")
    } else {
        info!("installing homebrew...");

        let resp = reqwest::blocking::get(
            "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh",
        )
        .expect("request failed");
        let body = resp.text().expect("body invalid");
        let mut out = File::create("homebrew_install.sh").expect("failed to create file");
        copy(&mut body.as_bytes(), &mut out).expect("failed to copy content");

        info!("homebrew install script downloaded");
        let install_script = read_to_string("homebrew_install.sh")
            .expect("Should have been able to read the file");

        Command::new("/bin/bash")
            .env("NONINTERACTIVE", "1")
            .arg("-c")
            .arg(install_script)
            .spawn()
            .expect("homebrew installation failed to start");

        match remove_file(Path::new("homebrew_install.sh")) {
            Ok(_) => info!("homebrew install script removed"),
            Err(e) => error!("error removing homebrew script, {:?}", e)
        };
    }
}

fn link_dotfiles() {
    let _links: Vec<(&str, &str)> = vec![
        ("", "")
    ];
    let source = Path::new("/Users/jacobshu/Documents/test/text.md");
    let target = Path::new("/Users/jacobshu/Documents/test/d/text.md");
    match symlink(source, target) {
        Ok(_) => info!("symlink created, {:?} => {:?}", source, target),
        Err(e) => error!("error creating symlink, {:?} => {:?} : {:?}", source, target, e)
    };
}

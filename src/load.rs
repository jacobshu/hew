use std::process::Command;

pub fn init() {
    println!("run load");
}

fn install_homebrew() {
    Command::new("NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
    .spawn()
    .expect("homebrew installation failed to start");
}
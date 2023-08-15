use crate::utils::which;
use clap::ArgMatches;
use log::{debug, error, info};
use std::fs::{read_to_string, remove_file, File, canonicalize, create_dir_all, remove_dir_all};
use std::io::copy;
use std::os::unix::fs::symlink;
use std::path::Path;
use std::process::Command;
use serde_derive::Deserialize;
use toml::from_str;
use dirs;

#[derive(Deserialize)]
struct Data {
    dotfiles: Vec<Symlink>,
}

#[derive(Deserialize)]
struct Symlink {
    source: String,
    target: String,
}

pub fn init(args: &ArgMatches) {
    info!(
        "beginning load::  update: {:?} | link: {:?}",
        args.get_flag("update"),
        args.get_flag("link")
    );
    let symlink_config = "../config/symlinks.toml";
    install_homebrew();
    link_dotfiles(symlink_config);
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
        let install_script =
            read_to_string("homebrew_install.sh").expect("Should have been able to read the file");

        Command::new("/bin/bash")
            .env("NONINTERACTIVE", "1")
            .arg("-c")
            .arg(install_script)
            .spawn()
            .expect("homebrew installation failed to start");

        match remove_file(Path::new("homebrew_install.sh")) {
            Ok(_) => info!("homebrew install script removed"),
            Err(e) => error!("error removing homebrew script, {:?}", e),
        };
    }
}

fn create_symlink(link: Symlink) -> Result<String, String> {
    let is_empty: Result<_, &str> = match (link.source:= "", link.target:= "") {
        (true, true) => { 
            Err("must provide source and target")
        },
        (true, false) => { Err("must provide source") },
        (false, true) => { Err("must provide target") },
        (false, false) => Ok(()),
    };
   
    if is_empty.is_err() { return Err(is_empty.unwrap_err().to_string()) }

    let source: canonicalize(Path::new(&link.source)).unwrap();
    let target: dirs::home_dir().unwrap().join(&link.target);
   
    // is_dir and is_file imply existence, symlinks will return true for these as well
    let source_status = (source.is_dir(), source.is_file());
    let target_status = (target.is_dir(), target.is_file());

    let status: Result<_, String> =  match (source_status, target_status) {
        ((false, false), (_, _)) => { Err(format!("target {:?} does not exist", target)) },
        ((true, true), (_, _)) => { Err("not possible, cannot be dir and file".to_string()) }
        ((_, _), (true, true)) => { Err("not possible, cannot be dir and file".to_string()) }
        ((true, false), (_, true)) => {
            Err(format!("{:?} and {:?} are of different types", source, target))
        },
        ((false, true), (true, _)) => {
            Err(format!("{:?} and {:?} are of different types", source, target))
        },
        
        // source & target are files: remove target
        ((false, true), (_, true)) => {
            remove_file(&target).expect("failed to remove target file");
            Ok(())
        },
        
        // source is file, target doesn't exist
        ((false, true), (false, false)) => { 
            let mut target_dir = dirs::home_dir().unwrap().join(&link.target);
            target_dir.pop();
            info!("target file does not exist: {:?}... ensuring directory path exists: {:?}", target, target_dir);
            create_dir_all(target_dir.clone()).expect("failed to create directory path");
            Ok(())
        }, 
       
        // source is dir, target is dir, remove target
        ((true, false), (true, _)) => { 
            if target != dirs::home_dir().unwrap() { 
                remove_dir_all(&target).expect("failed to remove target directory");
            }
            Ok(())
        },

        // source is dir, target doesn't exist
        ((true, false), (false, false)) => { 
            let target_dir = dirs::home_dir().unwrap().join(&link.target);
            info!("target directory does not exist. creating it...");
            create_dir_all(target_dir).expect("failed to create target directory");
            Ok(())
        }, 
    };

    if status.is_err() { return Err(status.err().unwrap()) }

    return match symlink(&source, &target) {
        Ok(_) => Ok(format!("symlink created, {:?} => {:?}", source, target)),
        Err(e) => Err(format!(
            "error creating symlink, {:?} => {:?} : {:?}",
            source, target, e
        )),
    };
}

fn link_dotfiles(config: &str) {
    let filename = Path::new(&config);
    let contents = read_to_string(filename).expect("error reading symlink config file");
    let data: Data = from_str(&contents).expect("error loading data from config file");

    for link in data.dotfiles {
        match create_symlink(link) {
            Ok(s) => s,
            Err(e) => e
       }; 
    }
}


// create_dir_all("Desktop/symlinks/source/symlink_dir_case").expect("failed target dir setup");
// create_dir_all("Desktop/symlinks/target/symlink_dir").expect("failed target dir setup");
// create_dir_all("Desktop/symlinks/source/symlink_not_there_dir").expect("failed target dir setup");
// File::create("Desktop/symlinks/source/file_exists_case.yaml").expect("error creating file");
// File::create("Desktop/symlinks/target/file_exists.yaml").expect("error creating file");
// File::create("Desktop/symlinks/target/file_not_exists_case.yaml").expect("error creating file");
// File::create("Desktop/symlinks/target/file_and_dir_not_exsits_case.yaml").expect("error creating file");

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs::{File, create_dir_all};

    fn setup() {
        create_dir_all("Desktop/symlinks/source/dir").expect("failed source dir setup");
        create_dir_all("Desktop/symlinks/target/dir").expect("failed target dir setup");
    }

    #[test]
    fn no_source_or_target_fails() {
        let link: Symlink = Symlink { source: "".to_string(), target: "".to_string() };
        let output: Result<String, String> = create_symlink(link);
        assert!(output.is_err());
    }

    #[test]
    fn no_source_fails() {
        let link: Symlink = Symlink { source: "".to_string(), target: "Desktop/symlinks/target/file.yaml".to_string() };
        let output: Result<String, String> = create_symlink(link);
        assert!(output.is_err());
    }

    #[test]
    fn no_target_fails() {
        let link: Symlink = Symlink { source: "Desktop/symlinks/source/file.yaml".to_string(), target: "".to_string() };
        let output: Result<String, String> = create_symlink(link);
        assert!(output.is_err());
    }

    #[test]
    fn nonexistent_source_fails() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/not_there_file.yaml".to_string(), 
            target: "Desktop/symlinks/target/file.yaml".to_string() 
        };
        let output: Result<String, String> = create_symlink(link);
        assert!(output.is_err());
    }
   
    #[test]
    fn file_to_dir_fails() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/file.yaml".to_string(),
            target: "Desktop/symlinks/target/dir".to_string(),
        };
    }

    #[test]
    fn dir_to_file_fails() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/dir".to_string(),
            target: "Desktop/symlinks/target/file.yaml".to_string(),
        };
    }

    #[test]
    fn file_exists_case() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/file_exists_case.yaml".to_string(),
            target: "Desktop/symlinks/target/file_exists.yaml".to_string(),
        };
    }
    
    #[test]
    fn target_file_noexistent() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/file_not_exists_case.yaml".to_string(),
            target: "Desktop/symlinks/target/not_there_file.yaml".to_string(),
        };
    }
    
    #[test]
    fn create_full_path() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/file_and_dir_not_exsits_case.yaml".to_string(),
            target: "Desktop/symlinks/target/not_there_dir/not_there_file.yaml".to_string(),
        };
    }
    
    #[test]
    fn dir_case() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/symlink_dir_case".to_string(),
            target: "Desktop/symlinks/target/symlink_dir".to_string(),
        };
    }

    #[test]
    fn nonexistent_dir_case() {
        let link: Symlink = Symlink { 
            source: "Desktop/symlinks/source/symlink_not_there_dir".to_string(),
            target: "Desktop/symlinks/target/not_there_dir".to_string(),
        };
    }
}

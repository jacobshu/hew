use std::env;
use std::path::Path;
pub mod which;

pub fn which(target: &str) -> bool {
    let env_paths = env::var("PATH").expect("couldn't get $PATH");
    let str_paths: Vec<String> = env_paths.split(":").map(|s| s.to_string()).collect();
    for path in str_paths.iter() {
        if Path::new(path).join(target).exists() {
            return true;
        }
    }
    return false;
}
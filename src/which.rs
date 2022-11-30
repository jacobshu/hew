use std::env;
use std::path::Path;

pub fn paths_in_path() -> Vec<String> {
    let env_paths = env::var("PATH").expect("couldn't get $PATH");
    let str_paths: Vec<String> = env_paths.split(":").map(|s| s.to_string()).collect();
    println!("{:?}", str_paths);
    return str_paths
}

fn find_in_path(target: &str) -> bool {
    let env_paths = paths_in_path();
    for path in env_paths.iter() {

    }
    let path = Path::new()
//    messages := make(chan string)
//go func() {
//        defer close(messages)
//var paths = EnviromentPaths()
//for _, target := range targets {
//    for _, path := range paths {
//        var target = path + "/" + target
//    if Exists(target) {
//        messages <- target
//    }
//    }
//}
//    }()
    return false
}
//
//func main() {
//    if len(os.Args) > 1 {
//        found := FindInPath(os.Args[1:])
//        for path := range found {
//            fmt.Println(path)
//        }
//    } else {
//        fmt.Println("usage: which program ...")
//    os.Exit(1)
//    }
//}
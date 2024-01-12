package main

import "github.com/jacobshu/hew/cmd"

//go:embed symlinks.toml
var symlinksToml string

//go:embed .env
var uri string


func main() {
  cmd.Execute()
}

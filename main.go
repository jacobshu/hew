package main

import (
  "embed"

  "hew.jacobshu.dev/cmd"
)

//go:embed symlinks.toml
var symlinksToml string

type ErrMsg struct{ err error }

func main() {
  cmd.Start()
}

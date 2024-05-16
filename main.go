package main

import (
  "hew.jacobshu.dev/cmd"
)

type ErrMsg struct{ err error }

func main() {
  cmd.Start()
}

package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"
)

func Start() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()
	log.SetOutput(f)

	if err := BuildCmdTree().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

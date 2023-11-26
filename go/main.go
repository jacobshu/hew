package main

import (
	"context"
  _ "embed"
	"fmt"
	"log"
	"os"
  "strings"
	"time"

  "github.com/jackc/pgx/v5"
)

//go:embed symlinks.toml
var symlinksToml string

//go:embed .env
var uri string

func openDB() *devDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
  uri = strings.ReplaceAll(uri, "\n", "")
	if uri == "" {
    log.Println("Warning:\n\tTo use the task manager, you must set your 'MONGODB_URI' environment variable.\n\tSee https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")

    client, err := pgx.Connect(ctx, uri)
    if err != nil {
      panic(err)
    } 

    mock := devDB{
      db:  client,
      ctx: ctx,
      closeDb: func() {
        cancel()
        if err := client.Close(ctx); err != nil {
          panic(err)
        }
      },
    }
    return &mock
	}


	client, err := pgx.Connect(ctx, uri)
	if err != nil {
		panic(err)
	}

	t := devDB{
		db:  client,
		ctx: ctx,
		closeDb: func() {
			cancel()
			if err := client.Close(ctx); err != nil {
				panic(err)
			}
		},
	}
	return &t
}

var devDb = openDB()

func main() {
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

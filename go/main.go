package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func openDB() *devDB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	t := devDB{
		db:  client,
		ctx: ctx,
		closeDb: func() {
			cancel()
			if err := client.Disconnect(ctx); err != nil {
				panic(err)
			}
		},
	}
	return &t
}

var devDb = openDB()

func main() {
  f, err := os.OpenFile("logfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

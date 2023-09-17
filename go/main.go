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
  defer cancel()

  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

  t := devDB{db: client, ctx: ctx}
  return &t
}

var Db = openDB()

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


package kinsta

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbName := strings.Join([]string{"file:", homeDir, ".config/hew/kinsta.db"}, "")

	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()

}

func Companies() {
  sites, err := Kinsta("GET", "/sites?company=fbd13128-664b-4cd3-9f1e-725a1a4d6f54", nil)
  if err != nil {
    fmt.Printf("error getting sites %v", err)
  }
  
  fmt.Printf("sites: \n%#v", sites)
}

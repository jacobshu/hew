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



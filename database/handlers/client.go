package handlers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Client struct {
	*sql.DB
}

/* The following function creates a new client object instance on which methods
 * are called to read data from the database.
 *
 * Returns: client (database) object
 */
func New() Client {
	connStr := os.Getenv("POSTGRE_CONN_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("[FATAL ERROR] Cannot open a connection to the database.")
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("[FATAL ERROR] Cannot communicate to the database.")
		panic(err)
	}

	return Client{db}
}

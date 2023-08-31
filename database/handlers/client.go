package handlers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func New() *sql.DB {

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

	return db
}

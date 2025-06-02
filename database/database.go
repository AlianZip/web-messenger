package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./messenger.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQLite")

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            hash TEXT NOT NULL,
			timestamp INT NOT NULL
        )
    `)
	if err != nil {
		log.Fatal(err)
	}
}

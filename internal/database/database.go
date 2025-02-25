package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	initTableQuery := `
	create table if not exists subscribers (
		id integer not null primary key autoincrement,
		email text,
		created_at timestamp default current_timestamp
	);`

	_, err = DB.Exec(initTableQuery)
	if err != nil {
		log.Fatalf("error creating table: %q: %s\n", err, initTableQuery)
	}
}

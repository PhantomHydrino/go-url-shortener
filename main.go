package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/PhantomHydrino/go-url-shortener/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:dummydb.sqlite?cache=shared")
	if err != nil {
		log.Println("Error is:", err)
		return
	}
	defer db.Close()

	// Create a table
	createTable := `
        CREATE TABLE IF NOT EXISTS urls (
            original_link TEXT,
			short_id TEXT
        )
    `
	if _, err = db.Exec(createTable); err != nil {
		log.Println("Error creating table:", err)
		return
	}

	// Create a new instance of the service
	s := utils.NewService(db)

	http.HandleFunc("/", s.ShortenURLHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

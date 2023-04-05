package main

import (
  "fmt"  
	"database/sql"
  _	"github.com/mattn/go-sqlite3"

)

func setupDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	createTableQuery := `
CREATE TABLE IF NOT EXISTS entries(
  id INTEGER PRIMARY KEY,
  timestamp TEXT NOT NULL,
  content TEXT NOT NULL
  )
  `

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func saveToDb(db *sql.DB, content [][]string) error {
	insertQuery := "INSERT INTO entries (timestamp, content) VALUES (datetime(\"now\", \"localtime\"), ?)"
	for _, entry := range content {
		_, err := db.Exec(insertQuery, entry[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func listDbContents(db *sql.DB) error {
	rows, err := db.Query("SELECT timestamp, content FROM entries")
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Printf("Timestamp\t\t\tTIL\n")
	var timestamp, content string
	for rows.Next() {
		err := rows.Scan(&timestamp, &content)
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", timestamp, content)
	}
	return rows.Err()
}



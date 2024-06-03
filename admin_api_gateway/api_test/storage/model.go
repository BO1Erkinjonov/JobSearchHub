package storage

import (
	"database/sql"
	"log"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Comment struct {
	ID          string `json:"ID"`
	Description string `json:"Description"`
	PstId       string `bson:"PstId"`
}

func ConnDB() *sql.DB {
	connStr := "postgres://postgres:123@localhost/testdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

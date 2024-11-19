package main

import (
	"database/sql"
)

func main() {
	db, err := sql.Open("mysql", "books:books@tcp(host:3306)/books")

	if err != nil {
		panic("Error opening db")
	}
}

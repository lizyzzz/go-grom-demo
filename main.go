package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Password string
}

func main() {
	db, err := sql.Open("mysql", "root:123456@(localhost:3306)/webserver")
	if err != nil {
		panic(err)
	}

	row := db.QueryRowContext(context.Background(), "SELECT username, password FROM user WHERE id = 1")
	if row.Err() != nil {
		panic(err)
	}

}

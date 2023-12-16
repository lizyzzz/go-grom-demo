package main

import (
	gormdemo "go-gorm-demo/gorm-demo"
	sqldemo "go-gorm-demo/sql-demo"
)

type User struct {
	Username string
	Password string
}

func main() {
	sqldemo.SqlDemo()
	gormdemo.CreateTable()
	gormdemo.Connect()
}

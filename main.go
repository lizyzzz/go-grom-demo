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
	// gormdemo.CreateTable()
	// gormdemo.InsertTable()
	// gormdemo.QueryTable()
	// gormdemo.UpdateTable()
	// gormdemo.DeleteTable()
	// gormdemo.TestHook()
	// gormdemo.PrepareData()
	// gormdemo.AdvancedQuery()
	gormdemo.One2MoreCreateTable()
	gormdemo.One2MoreInsertTable()
	// gormdemo.One2MoreQuery()
	// gormdemo.One2MoreDelete()
}

package sqldemo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Password string
}

func SqlDemo() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/webserver")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 返回值只有一行的查询
	id := 1
	row := db.QueryRowContext(context.Background(), "SELECT username, password FROM user WHERE id = ?;", id)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			fmt.Println(row.Err())
		} else {
			panic(row.Err())
		}
	}

	// 解析结果
	var user User
	if err = row.Scan(&user.Username, &user.Password); err != nil {
		panic(err)
	}
	fmt.Println(user)

	// 返回值有多行的查询
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM user;")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
		} else {
			panic(err)
		}
	}

	for rows.Next() {
		var ID int32
		var name string
		var pwd string
		if err := rows.Scan(&ID, &name, &pwd); err != nil {
			panic(err)
		}
		fmt.Println(ID, name, pwd)
	}
	if rows.Err() != nil {
		panic(rows.Err())
	}

	// 支持事务
}

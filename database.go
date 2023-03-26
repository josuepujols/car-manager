package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDataBase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost)/crud_go")
	if err != nil {
		panic(err.Error())
	}
	return db, err
}

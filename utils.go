package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDb() (*sql.DB, error) {
	return sql.Open("mysql", "root:root@tcp(db:3306)/blobber?parseTime=true&charset=utf8") //fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname))
}

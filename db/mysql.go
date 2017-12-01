package db

import (
	 _ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

func GetDb(user, password, dbname string) (*sql.DB, error){
	datasourceName := fmt.Sprintf("%s:%s@/%s",user, password, dbname)
	return sql.Open("mysql", datasourceName)
}
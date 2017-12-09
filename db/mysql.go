package db

import (
	 _ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	umy "upper.io/db.v3/mysql"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetDb(user, password, dbname string) (*sql.DB, error) {
	datasourceName := fmt.Sprintf("%s:%s@/%s",user, password, dbname)
	return sql.Open("mysql", datasourceName)
}

func GetUpperDB(user, password, host, dbname string) (sqlbuilder.Database, error) {
	var settings = umy.ConnectionURL{
		User:     user,
		Password: password,
		Database: dbname,
		Host: host,
	}

	sess, err := umy.Open(settings)
	//sess.SetLogging(true)


	return sess, err
}
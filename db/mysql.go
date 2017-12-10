package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"upper.io/db.v3/lib/sqlbuilder"
	umy "upper.io/db.v3/mysql"
)

func GetDb(user, password, dbname string) (*sql.DB, error) {
	datasourceName := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	return sql.Open("mysql", datasourceName)
}

func GetUpperDB(user, password, host, dbname string) (sqlbuilder.Database, error) {
	var settings = umy.ConnectionURL{
		User:     user,
		Password: password,
		Database: dbname,
		Host:     host,
	}

	sess, err := umy.Open(settings)
	//sess.SetLogging(true)

	return sess, err
}

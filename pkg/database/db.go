package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var conn *sqlx.DB

func InitDB() error {
	db, err := sqlx.Open("postgres", "user=whatsub dbname=whatsub password=whatsub sslmode=disable")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	conn = db
	return nil
}

func GetConn() sqlx.DB {
	return *conn
}

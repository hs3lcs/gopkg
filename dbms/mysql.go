package dbms

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DbConn struct {
	DB *sql.DB
	Tx *sql.Tx
}

func Connect(mode string) (*sql.DB, error) {
	var host, user, pass string
	if mode == "RW" {
		host = Config.MYSQL_RW_HOST
		user = Config.MYSQL_RW_USER
		pass = Config.MYSQL_RW_PASS
	} else {
		host = Config.MYSQL_RO_HOST
		user = Config.MYSQL_RO_USER
		pass = Config.MYSQL_RO_PASS
	}
	cfg := mysql.Config{
		Net:    "tcp",
		Addr:   host,
		User:   user,
		Passwd: pass,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return db, err
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, err
}

func Transaction() (DbConn, error) {
	var conn DbConn
	var err error
	conn.DB, err = Connect("RW")
	if err != nil {
		return conn, err
	}
	conn.Tx, err = conn.DB.Begin()
	return conn, err
}

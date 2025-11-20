// MySQLの接続管理部分
package db

import (
	"database/sql"
	"fmt"

	 _ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Conn *sql.DB
}

func NewMySQLConnection(user, pass, host, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// コネクション確立チェック
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &Database{
		Conn: conn,
	}, nil
}
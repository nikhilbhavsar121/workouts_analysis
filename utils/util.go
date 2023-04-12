package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func GetDBConnection() (conn *sql.DB, err error) {

	dbport := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	passwd := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")

	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s",
		user, passwd, dbport, databaseName))
	if err != nil {
		fmt.Println("Cannot connect to database", err)
		return
	}
	return conn, err
}

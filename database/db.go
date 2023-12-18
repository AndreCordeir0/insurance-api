package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetConnection() error {
	fmt.Println("Tentando conexao...")

	var err error

	config := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Addr:   os.Getenv("MYSQL_HOST"),
		DBName: os.Getenv("MYSQL_DATABASE"),
		Net:    "tcp",
	}
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}
	return nil
}

func GetDb() *sql.DB {
	return db
}

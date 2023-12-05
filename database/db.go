package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AndreCordeir0/insurance-api/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetConnection() error {
	fmt.Println("Tentando conexao...")
	config := &config.ConfigDB{
		DB_USER:     os.Getenv("MYSQL_USER"),
		DB_PASSWORD: os.Getenv("MYSQL_PASSWORD"),
		DB_HOST:     os.Getenv("MYSQL_HOST"),
		DB_PORT:     "3306",
		DB_NAME:     os.Getenv("MYSQL_DATABASE"),
	}
	var err error
	configConnection := config.GetDataSorceConnectionName()
	println(configConnection)
	db, err = sql.Open("mysql", configConnection)
	if err != nil {
		return err
	}
	return nil
}

func GetDb() *sql.DB {
	return db
}

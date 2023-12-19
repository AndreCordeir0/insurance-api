package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func GetConnection() *sql.DB {
	fmt.Println("Tentando conexao...")

	config := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Addr:   os.Getenv("MYSQL_HOST"),
		DBName: os.Getenv("MYSQL_DATABASE"),
		Net:    "tcp",
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal("Erro ao conectar no banco de dados", err)
		panic(err)
	}
	return db
}

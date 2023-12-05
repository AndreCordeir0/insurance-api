package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AndreCordeir0/insurance-api/database"
	"github.com/AndreCordeir0/insurance-api/models"
)

func handler(w http.ResponseWriter, r *http.Request) {
	errConnection := database.GetConnection()
	if errConnection != nil {
		log.Fatal("Erro ao conectar no banco de dados", errConnection)
		panic(errConnection)
	}
	database.GetDb()
	transaction, err := database.GetDb().Begin()
	if err != nil {
		log.Fatal(err)
	}
	insurance := &models.Insurance{}
	result, err := transaction.Exec("INSERT INTO TB_INSURANCE (age, dependents, income, marital_status) VALUES (?, ?, ?, ?)",
		insurance.Age, insurance.Dependents, insurance.Income, insurance.MaritalStatus)
	if err != nil {
		transaction.Rollback()
		log.Fatal(err)
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		transaction.Rollback()
		log.Fatal(errCommit)
	}
	database.GetDb().Close()
	println("Resultado", result)
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Escutando na porta :8080...")
	err := http.ListenAndServe(":8080", nil)
	fmt.Println("aaaaa")
	if err != nil {
		log.Fatal("Erro")
		panic(err)
	}

}

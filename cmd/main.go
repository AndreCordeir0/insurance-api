package main

import (
	"fmt"
	"log"

	"github.com/AndreCordeir0/insurance-api/database"
	"github.com/AndreCordeir0/insurance-api/models"
	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	errConnection := database.GetConnection()
	if errConnection != nil {
		log.Fatal("Erro ao conectar no banco de dados", errConnection)
		panic(errConnection)
	}
	transaction, err := database.GetDb().Begin()
	defer database.GetDb().Close()

	if err != nil {
		log.Fatal(err)
	}
	insurance := &models.Insurance{
		Age:           25,
		Dependents:    2,
		Income:        0,
		MaritalStatus: "not married",
	}
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
	println("Resultado", result)
	c.JSON(200, result)
}

func main() {
	r := gin.Default()
	r.GET("/teste", handler)
	fmt.Println("Escutando na porta :8080...")
	r.Run(":8080")
}

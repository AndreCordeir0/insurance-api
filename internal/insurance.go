package internal

import (
	"log"

	"github.com/AndreCordeir0/insurance-api/database"
	"github.com/gin-gonic/gin"
)

type Insurance struct {
	ID            int    `json:"id"`
	Age           int    `json:"age"`
	Income        int    `json:"income"`
	MaritalStatus string `json:"marital_status"`
	Dependents    int    `json:"dependents"`
}

func Insert(c *gin.Context) {
	//TODO - Criar um serviço para inserir no banco de dados
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

	insurance := &Insurance{
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

func GetAll(c *gin.Context) {
	errConnection := database.GetConnection()
	if errConnection != nil {
		log.Fatal("Erro ao conectar no banco de dados", errConnection)
		panic(errConnection)
	}
	rows, err := database.GetDb().Query("SELECT * FROM TB_INSURANCE")
	defer rows.Close()
	if err != nil {
		c.JSON(500, err)
	}
	var insurances []Insurance
	for rows.Next() {
		var ins Insurance
		if err := rows.Scan(&ins.ID, &ins.Age, &ins.Dependents, &ins.Income, &ins.MaritalStatus); err != nil {
			c.JSON(500, err)
		}
		insurances = append(insurances, ins)
	}
	c.JSON(200, insurances)
}

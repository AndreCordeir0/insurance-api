package internal

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/AndreCordeir0/insurance-api/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	house_status   []string = []string{"mortgaged", "owned"}
	marital_status []string = []string{"married", "single"}
)

// Unica forma que encontrei para contornar : 0 -> null ao usar required
type Insurance struct {
	ID            int     `json:"id"`
	Age           *int    `json:"age" validate:"required,min=0"`
	Income        *int    `json:"income" validate:"required,min=0"`
	MaritalStatus string  `json:"marital_status" validate:"required"`
	Dependents    *int    `json:"dependents" validate:"required,min=0"`
	RiskQuestion  []int   `json:"risk_question" validate:"required"`
	Vehicle       Vehicle `json:"vehicle" validate:"required"`
	House         House   `json:"house" validate:"required"`
}

type Vehicle struct {
	Year *int `json:"year" validate:"required"`
}

type House struct {
	OwnershipStatus string `json:"ownership_status" validate:"required"`
}

func Insert(c *gin.Context) {
	databaseConnection := database.GetConnection()
	transaction, err := databaseConnection.Begin()
	defer databaseConnection.Close()
	if err != nil {
		log.Fatal(err)
	}

	var insurance Insurance

	if err := c.ShouldBindJSON(&insurance); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	validateError := validateInsurance(&insurance)
	if validateError != nil {
		c.JSON(400, gin.H{"error": validateError.Error()})
		return
	}
	// idVehicle := insertVehicle()
	// idHouse := insertHouse()
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
	id, _ := result.LastInsertId()
	c.JSON(200, id)
}

func insertVehicle() int64 {

}

func insertHouse() int64 {

}

func validateInsurance(insurance *Insurance) error {
	validator := validator.New()
	err := validator.Struct(insurance)
	if err != nil {
		fmt.Println("Erro:", err)
		return err
	}
	index := sort.SearchStrings(house_status, insurance.House.OwnershipStatus)
	indexMarital := sort.SearchStrings(marital_status, insurance.MaritalStatus)

	boole := (index < len(house_status) && house_status[index] == insurance.House.OwnershipStatus && index < len(marital_status) && marital_status[indexMarital] == insurance.MaritalStatus)
	if !boole {
		return errors.New("invalid ownership status for the house or marital status")
	}
	return nil
}

func GetAll(c *gin.Context) {
	database := database.GetConnection()
	rows, err := database.Query("SELECT * FROM TB_INSURANCE")
	if err != nil {
		c.JSON(500, err)
		return
	}
	defer rows.Close()
	var insurances []Insurance
	for rows.Next() {
		var ins Insurance
		if err := rows.Scan(&ins.ID, &ins.Age, &ins.Dependents, &ins.Income, &ins.MaritalStatus); err != nil {
			c.JSON(500, err)
			return
		}
		insurances = append(insurances, ins)
	}
	c.JSON(200, insurances)
}

package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/AndreCordeir0/insurance-api/database"
	"github.com/AndreCordeir0/insurance-api/utils"
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
	IdVehicle     int     `json:"-"`
	IdHouse       int     `json:"-"`
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
	//Pequeno truque para rollback, se o commit já houver sido confirmado o rollback não é feito (:
	defer transaction.Rollback()
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
	idVehicle, vehicleError := insertVehicle(&insurance, transaction)
	idHouse, houseError := insertHouse(&insurance, transaction)
	if vehicleError != nil {
		c.JSON(400, gin.H{"error": vehicleError.Error()})
		return
	}
	if houseError != nil {
		c.JSON(400, gin.H{"error": houseError.Error()})
		return
	}

	result, err := transaction.Exec("INSERT INTO TB_INSURANCE (age, dependents, income, marital_status, id_vehicle, id_house) VALUES (?, ?, ?, ?, ?, ?)",
		insurance.Age, insurance.Dependents, insurance.Income, insurance.MaritalStatus, idVehicle, idHouse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	errCommit := transaction.Commit()
	if errCommit != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	id, _ := result.LastInsertId()
	c.JSON(200, id)
}

func insertVehicle(insurance *Insurance, transaction *sql.Tx) (int64, error) {
	println(*insurance.Vehicle.Year)
	id, err := utils.AbstractInsert[int]("TB_VEHICLE", []string{"year"}, []int{*insurance.Vehicle.Year}, transaction)
	return id, err
}

func insertHouse(insurance *Insurance, transaction *sql.Tx) (int64, error) {
	id, err := utils.AbstractInsert[string]("TB_HOUSE", []string{"ownership_status"}, []string{insurance.House.OwnershipStatus}, transaction)
	return id, err
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
		c.JSON(500, err.Error())
		return
	}
	defer rows.Close()
	var insurances []Insurance
	for rows.Next() {
		var ins Insurance
		if err := rows.Scan(&ins.ID, &ins.Age, &ins.Dependents, &ins.Income, &ins.MaritalStatus, &ins.IdVehicle, &ins.IdHouse); err != nil {
			c.JSON(500, err.Error())
			return
		}
		insurances = append(insurances, ins)
	}
	c.JSON(200, insurances)
}

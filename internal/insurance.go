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

const (
	Economic    = "economic"
	Regular     = "regular"
	Responsible = "responsible"
	Ineligible  = "ineligible"
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

type RiskScore struct {
	Auto       string `json:"auto"`
	Disability string `json:"disability"`
	Home       string `json:"home"`
	Life       string `json:"life"`
}

type RiskScoreNumber struct {
	Auto       int
	Disability int
	Home       int
	Life       int
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
func CalculateScore(insurance *Insurance) {
	riskSum := Sum(insurance.RiskQuestion)
	var riskScoreNumber *RiskScoreNumber = &RiskScoreNumber{
		Auto:       riskSum,
		Disability: riskSum,
		Home:       riskSum,
		Life:       riskSum,
	}
	//TODO
	// If the user has dependents, add 1 risk point to both the disability and life scores.
	// If the user is married, add 1 risk point to the life score and remove 1 risk point from disability.
	// If the user's vehicle was produced in the last 5 years, add 1 risk point to that vehicle’s score.
	var riskScore *RiskScore = &RiskScore{}
	DetermineInsuranceEligibility(insurance, riskScore)
	DetermineAgeEligibility(insurance, riskScore, riskScoreNumber)
	DetermineIncomeEligibility(insurance, riskScoreNumber)
	DetermineHouseEligibility(insurance, riskScoreNumber)
	DetermineIncomeEligibility(insurance, riskScoreNumber)
	DetermineIncomeEligibility(insurance, riskScoreNumber)

}

func DetermineInsuranceEligibility(insurance *Insurance, riskScore *RiskScore) {
	// Implemented - If the user doesn’t have income, vehicles or houses, she is ineligible for disability, auto, and home insurance, respectively.
	if *insurance.Income == 0 {
		riskScore.Disability = Ineligible
	}
	if insurance.Vehicle == (Vehicle{}) {
		riskScore.Auto = Ineligible
	}
	if insurance.House == (House{}) {
		riskScore.Home = Ineligible
	}
}

func DetermineAgeEligibility(insurance *Insurance, riskScore *RiskScore, riskScoreNumber *RiskScoreNumber) {
	age := *insurance.Age
	// Implemented - If the user is over 60 years old, she is ineligible for disability and life insurance.
	if age > 60 {
		riskScore.Disability = Ineligible
		riskScore.Life = Ineligible
	}

	// Implemented - If the user is under 30 years old, deduct 2 risk points from all lines of insurance. If she is between 30 and 40 years old, deduct 1.
	if age < 30 {
		riskScoreNumber.Auto = riskScoreNumber.Auto - 2
		riskScoreNumber.Disability = riskScoreNumber.Disability - 2
		riskScoreNumber.Home = riskScoreNumber.Home - 2
		riskScoreNumber.Life = riskScoreNumber.Life - 2
	} else if age >= 30 && age <= 40 {
		riskScoreNumber.Auto = riskScoreNumber.Auto - 1
		riskScoreNumber.Disability = riskScoreNumber.Disability - 1
		riskScoreNumber.Home = riskScoreNumber.Home - 1
		riskScoreNumber.Life = riskScoreNumber.Life - 1
	}
}
func DetermineIncomeEligibility(insurance *Insurance, riskScoreNumber *RiskScoreNumber) {
	// Implemented - If her income is above $200k, deduct 1 risk point from all lines of insurance.
	income := *insurance.Income
	if income > 200000 {
		riskScoreNumber.Auto = riskScoreNumber.Auto - 1
		riskScoreNumber.Disability = riskScoreNumber.Disability - 1
		riskScoreNumber.Home = riskScoreNumber.Home - 1
		riskScoreNumber.Life = riskScoreNumber.Life - 1
	}
}

func DetermineHouseEligibility(insurance *Insurance, riskScoreNumber *RiskScoreNumber) {
	// If the user's house is mortgaged, add 1 risk point to her home score and add 1 risk point to her disability score.
	isMortgated := house_status[0]
	if (House{}) != insurance.House && insurance.House.OwnershipStatus == isMortgated {
		riskScoreNumber.Home = riskScoreNumber.Home - 1
		riskScoreNumber.Disability = riskScoreNumber.Home - 1
	}
}

//TODO func deduct all income points

func Sum(array []int) int {
	var sum int
	for _, item := range array {
		sum += item
	}
	return sum
}

func GetPontuation(score int8) string {
	switch score {
	case 0:
		return Economic
	case 1, 2:
		return Regular
	default:
		return Responsible
	}
}

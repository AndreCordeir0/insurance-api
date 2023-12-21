package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInsurance_ValidInput(t *testing.T) {
	validInsurance := &Insurance{
		Age:           IntPointer(25),
		Income:        IntPointer(50000),
		MaritalStatus: "single",
		Dependents:    IntPointer(1),
		RiskQuestion:  []int{0, 1, 0},
		Vehicle: Vehicle{
			Year: IntPointer(2020),
		},
		House: House{
			OwnershipStatus: "owned",
		},
	}

	err := ValidateInsurance(validInsurance)

	assert.Nil(t, err, "Expected no error for valid input")
}
func TestValidateInsurance_InvalidHouseOwnershipStatus(t *testing.T) {
	invalidInsurance := &Insurance{
		Age:           IntPointer(25),
		Income:        IntPointer(50000),
		MaritalStatus: "single",
		Dependents:    IntPointer(1),
		RiskQuestion:  []int{0, 1, 0},
		Vehicle: Vehicle{
			Year: IntPointer(2020),
		},
		House: House{
			OwnershipStatus: "InvalidStatus",
		},
	}

	err := ValidateInsurance(invalidInsurance)

	assert.NotNil(t, err, "Expected error for invalid house ownership status")
	assert.Contains(t, err.Error(), "invalid ownership status for the house or marital status")
}

func TestValidateInsurance_InvalidRisk(t *testing.T) {
	insurance := ReturnValidInsurance()
	insurance.RiskQuestion = []int{1, 3, 0}

	err := ValidateInsurance(insurance)
	assert.NotNil(t, err, "Expected error for invalid house ownership status")
	assert.Contains(t, err.Error(), "invalid value in risk_question")
}

func ReturnValidInsurance() (insurance *Insurance) {
	insurance = &Insurance{
		Age:           IntPointer(25),
		Income:        IntPointer(50000),
		MaritalStatus: "single",
		Dependents:    IntPointer(1),
		RiskQuestion:  []int{0, 1, 0},
		Vehicle: Vehicle{
			Year: IntPointer(2020),
		},
		House: House{
			OwnershipStatus: "owned",
		},
	}
	return insurance
}

func IntPointer(i int) *int {
	return &i
}

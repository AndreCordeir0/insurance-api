package models

type Insurance struct {
	ID            int    `json:"id"`
	Age           int    `json:"age"`
	Income        int    `json:"income"`
	MaritalStatus string `json:"marital_status"`
	Dependents    int    `json:"dependents"`
}

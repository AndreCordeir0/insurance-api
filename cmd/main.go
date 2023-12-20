package main

import (
	"github.com/AndreCordeir0/insurance-api/internal"
	"github.com/gin-gonic/gin"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	r := gin.Default()
	r.POST("/insurance-risk", internal.InsuranceRiskEstimated)
	r.Run(":8080")
}

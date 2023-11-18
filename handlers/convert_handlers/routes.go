package convert_handlers

import (
	"github.com/attributeerror/currency-rates-service/database"
	"github.com/gin-gonic/gin"
)

func InitialiseRoutes(e *gin.Engine, db database.Database) {
	// app endpoints
	rootGroup := e.Group("/currency-rates-service")
	v1Group := rootGroup.Group("/v1")
	v1Group.GET("convert", GetConvertCurrencyFromBase(db))

	// liveness/readiness endpoint
	e.GET("liveness", LivenessProbe())
}

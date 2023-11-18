package convert_handlers

import (
	"github.com/attributeerror/currency-rates-service/database"
	"github.com/gin-gonic/gin"
)

func InitialiseRoutes(e *gin.Engine, db database.Database) {
	e.GET("convert", GetConvertCurrencyFromBase(db))
}

package convert_handlers

import (
	"github.com/attributeerror/currency-rates-service/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

func InitialiseRoutes(e *gin.Engine, db database.Database, sfGroup *singleflight.Group) {
	e.GET("convert", GetConvertCurrencyFromBase(db, sfGroup))
	e.GET("liveness", LivenessProbe())
}

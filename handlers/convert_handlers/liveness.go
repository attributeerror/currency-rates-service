package convert_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var LivenessProbe = func() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

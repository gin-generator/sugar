package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterApi(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

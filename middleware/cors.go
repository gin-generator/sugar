package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors
/**
 * @description: Cors handles Cross-Origin Resource Sharing (CORS) settings.
 */
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", getAllowOrigin())
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func getAllowOrigin() (origin string) {
	cors := []string{
		"127.0.0.1",
	}
	if len(cors) > 0 {
		origin = strings.Join(cors, ",")
	} else {
		origin = "*"
	}
	return origin
}

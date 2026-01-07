package bootstrap

import (
	"fmt"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-gonic/gin"
)

// RegisterRouter
/**
 * @description: router registration function type
 */
type RegisterRouter func(*gin.Engine)

// Http
/**
 * @description: http server struct
 */
type Http struct {
	*gin.Engine
}

// newHttp
/**
 * @description: create a new http server instance
 * @param {string} env
 * @return {*Http}
 */
func newHttp(env string) *Http {
	gin.SetMode(env)
	return &Http{
		Engine: gin.New(),
	}
}

// Run starts the HTTP server
func (h *Http) Run(app *foundation.Application) {
	cfg := app.Config

	name := cfg.App.Name
	host := cfg.App.Host
	port := cfg.App.Port

	fmt.Printf("%s serve start: %s:%d...\n", name, host, port)
	err := h.Engine.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic("Unable to start server, error: " + err.Error())
	}
}

// Use add middleware
/**
 * @description: add middleware to the http server
 * @param {...gin.HandlerFunc} middleware
 */
func (h *Http) Use(middleware ...gin.HandlerFunc) {
	h.Engine.Use(middleware...)
}

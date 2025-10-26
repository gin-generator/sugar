package bootstrap

import (
	"fmt"
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
 * @param {Mode} env
 * @return {*Http}
 */
func newHttp(env Mode) *Http {
	gin.SetMode(string(env))
	return &Http{
		Engine: gin.New(),
	}
}

// run
/**
 * @description: run http server
 * @param {Config} cfg
 */
func (h *Http) run(cfg *Config) {
	fmt.Println(fmt.Sprintf("%s serve start: %s:%d...",
		cfg.App.Name, cfg.App.Host, cfg.App.Port))
	err := h.Engine.Run(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	if err != nil {
		panic("Unable to start server, error: " + err.Error())
	}
}

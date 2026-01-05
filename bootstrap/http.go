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

// Run 实现 Server 接口
/**
 * @description: run http server
 * @param {Config} cfg
 */
func (h *Http) Run(cfg *Config) {
	fmt.Println(fmt.Sprintf("%s serve start: %s:%d...",
		cfg.App.Name, cfg.App.Host, cfg.App.Port))
	err := h.Engine.Run(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	if err != nil {
		panic("Unable to start server, error: " + err.Error())
	}
}

// GetEngine 实现 Server 接口
/**
 * @description: get gin engine
 * @return {*gin.Engine}
 */
func (h *Http) GetEngine() *gin.Engine {
	return h.Engine
}

// Use 实现 Server 接口
/**
 * @description: use middleware
 * @param {...gin.HandlerFunc} middleware
 */
func (h *Http) Use(middleware ...gin.HandlerFunc) {
	h.Engine.Use(middleware...)
}

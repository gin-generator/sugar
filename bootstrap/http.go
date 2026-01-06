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

// Run 实现 Server 接口
/**
 * @description: 运行 HTTP 服务器
 * @param {*foundation.Application} app
 */
func (h *Http) Run(app *foundation.Application) {
	cfg, _ := app.GetConfig("app")
	appCfg := cfg.(map[string]interface{})

	name := appCfg["name"].(string)
	host := appCfg["host"].(string)
	port := appCfg["port"].(int)

	fmt.Printf("%s serve start: %s:%d...\n", name, host, port)
	err := h.Engine.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic("Unable to start server, error: " + err.Error())
	}
}

// Use 添加中间件
/**
 * @description: 使用中间件
 * @param {...gin.HandlerFunc} middleware
 */
func (h *Http) Use(middleware ...gin.HandlerFunc) {
	h.Engine.Use(middleware...)
}

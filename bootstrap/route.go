package bootstrap

import (
	"github.com/gin-generator/sugar/app/demo/route"
	"github.com/gin-gonic/gin"
)

func RegisterDemoApiRoute(e *gin.Engine) {
	route.RegisterApi(e)
}

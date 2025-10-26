package bootstrap

import (
	"github.com/gin-generator/sugar/app/demo/route"
	"github.com/gin-gonic/gin"
)

func RegisterApiRoute(e *gin.Engine) {
	route.RegisterApi(e)
}

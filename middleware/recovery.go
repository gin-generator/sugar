package middleware

import (
	"errors"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// Recovery
/**
 * @Description: use zap.Error to record Panic and call stack
 */
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				// connection was aborted, we can't write a status to the client.
				var brokenPipe bool
				var ne *net.OpError
				if errors.As(err.(error), &ne) {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// connection was aborted, we can't write a status to the client.
				if brokenPipe {
					logger.Log.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				// if the connection is not aborted, we can log the panic stacktrace
				logger.Log.Error("recovery from panic",
					zap.Time("time", time.Now()),               // 记录时间
					zap.Any("error", err),                      // 记录错误信息
					zap.String("request", string(httpRequest)), // 请求
					zap.Stack("stacktrace"),                    // 记录调用堆栈信息
				)

				// 返回 500000 状态码
				// http.Alert500(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			}
		}()
		c.Next()
	}
}

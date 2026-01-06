package middleware

import (
	"errors"
	"github.com/gin-generator/sugar/services/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// Recovery uses zap.Error to record Panic and call stack
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
					zap.Time("time", time.Now()),               // record time
					zap.Any("error", err),                      // record error
					zap.String("request", string(httpRequest)), // request
					zap.Stack("stacktrace"),                    // record stack trace
				)

				// return 500 status code
				// http.Alert500(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			}
		}()
		c.Next()
	}
}

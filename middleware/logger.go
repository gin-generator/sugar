package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-generator/sugar/package/logger"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger
/**
 * @Description: record request log
 */
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if logger.Log.Log.Level() > zap.ErrorLevel {
			c.Next()
			return
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		start := time.Now()
		c.Next()

		cost := time.Since(start)
		responseStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", responseStatus),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", microsecondsStr(cost)),
		}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			logFields = append(logFields, zap.String("Request Body", string(requestBody)))
			logFields = append(logFields, zap.String("Response Body", w.body.String()))
		}

		if responseStatus > http.StatusBadRequest && responseStatus <= http.StatusUnavailableForLegalReasons {
			logger.Log.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
		} else if responseStatus >= http.StatusInternalServerError && responseStatus <= http.StatusNetworkAuthenticationRequired {
			logger.Log.Error("HTTP Error "+cast.ToString(responseStatus), logFields...)
		} else {
			if c.Request.MultipartForm == nil {
				logger.Log.Debug("HTTP Access Log", logFields...)
			}
		}
	}
}

func microsecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

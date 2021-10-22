package ginext

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/praslar/cloud0/common"
	"gitlab.com/praslar/cloud0/logger"
)

func AccessLogMiddleware(env string) gin.HandlerFunc {
	l := logger.WithField("env", env)
	extractHeaders := []string{"x-forwarded-for", common.HeaderTenantID, common.HeaderUserID, common.HeaderXRequestID}

	return func(c *gin.Context) {

		if c.Request.URL.Path == "/status" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		defer func() {
			latency := time.Since(start).Milliseconds()
			l = l.
				WithField("status", c.Writer.Status()).
				WithField("method", c.Request.Method).
				WithField("path", path).
				WithField("ip", c.ClientIP()).
				WithField("latency", latency).
				WithField("user-agent", c.Request.UserAgent())

			for _, header := range extractHeaders {
				if v := c.GetHeader(header); v != "" {
					l = l.WithField(header, v)
				}
			}

			l.Infof("acesss log")
		}()

		c.Next()
	}
}

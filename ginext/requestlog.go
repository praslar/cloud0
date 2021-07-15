package ginext

import (
	"time"

	"github.com/gin-gonic/gin"
	. "gitlab.com/goxp/cloud0/common"
	"gitlab.com/goxp/cloud0/logger"
)

func AccessLogMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	defer func() {
		latency := time.Since(start).Milliseconds()
		l := logger.
			WithField("status", c.Writer.Status()).
			WithField("method", c.Request.Method).
			WithField("path", path).
			WithField("ip", c.ClientIP()).
			WithField("latency", latency).
			WithField("user-agent", c.Request.UserAgent()).
			WithField("x-request-id", c.GetString(HeaderXRequestID)).
			WithField("x-user-id", c.GetString(HeaderUserID))

		if v := c.GetHeader("X-Forwarded-For"); v != "" {
			l = l.WithField("x-forwarded-for", v)
		}
		l.Infof("acesss log")
	}()

	c.Next()
}

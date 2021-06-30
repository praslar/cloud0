package ginext

import (
	"time"

	"github.com/gin-gonic/gin"
	. "gitlab.com/goxp/cloud0/common"
	"gitlab.com/goxp/cloud0/log"
)

func RequestLogMiddleware(tag string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logRequest(ctx, tag)
	}
}

func logRequest(c *gin.Context, tag string) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	defer func() {
		go func() {
			latency := time.Since(start).Milliseconds()
			l := log.Tag(tag).
				WithField("status", c.Writer.Status()).
				WithField("method", c.Request.Method).
				WithField("path", path).
				WithField("ip", c.ClientIP()).
				WithField("latency", latency).
				WithField("user-agent", c.Request.UserAgent()).
				WithField("x-request-id", c.GetString(HeaderXRequestID)).
				WithField("proto", c.Request.Proto).
				WithField("x-user-id", c.GetInt64(HeaderUserID))

			if v := c.GetHeader("X-Forward-For"); v != "" {
				l = l.WithField("x-forward-for", v)
			}
			if v := c.GetString("upstream"); v != "" {
				l = l.WithField("upstream", v)
			}
			l.Infof("acesss log")
		}()
	}()

	c.Next()
}
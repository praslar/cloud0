package ginext

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "gitlab.com/goxp/cloud0/common"
)

// AuthRequiredMiddleware is required the request has to have x-user-id in header
// (it usually set by API Gateway)
func AuthRequiredMiddleware(c *gin.Context) {
	headers := struct {
		UserID   int64  `header:"x-user-id" validate:"required,min=1"`
		UserMeta string `header:"x-user-meta"`
	}{}
	if c.ShouldBindHeader(&headers) != nil {
		_ = c.Error(NewError(http.StatusUnauthorized, "unauthorized"))
		c.Abort()
		return
	}

	c.Set(HeaderUserID, headers.UserID)
	c.Set(HeaderUserMeta, headers.UserMeta)

	c.Next()
}

type Int64Getter interface {
	GetInt64(key string) int64
}

// GetUserID returns the user ID embedded in Gin context
func GetUserID(c Int64Getter) int64 {
	return c.GetInt64(HeaderUserID)
}

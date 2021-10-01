package ginext

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	. "gitlab.com/goxp/cloud0/common"
)

// AuthRequiredMiddleware is required the request has to have x-user-id in header
// (it's usually set by API Gateway)
func AuthRequiredMiddleware(c *gin.Context) {
	headers := struct {
		UserID   string `header:"x-user-id" validate:"required,min=1"`
		UserMeta string `header:"x-user-meta"`
		TenantID uint64 `header:"x-tenant-id"`
	}{}
	if c.ShouldBindHeader(&headers) != nil {
		_ = c.Error(NewError(http.StatusUnauthorized, "unauthorized"))
		c.Status(http.StatusUnauthorized) // in case of we don't use this middleware with ErrorHandler
		c.Abort()
		return
	}

	c.Set(HeaderUserID, headers.UserID)
	c.Set(HeaderUserMeta, headers.UserMeta)
	c.Set(HeaderTenantID, headers.TenantID)

	c.Next()
}

type GetStringer interface {
	GetString(key string) string
}

// GetUserID returns the user ID embedded in Gin context
func GetUserID(c GetStringer) string {
	return c.GetString(HeaderUserID)
}

func Uint64HeaderValue(c *gin.Context, headerName string) uint64 {
	sValue := c.GetHeader(headerName)
	if sValue == "" {
		return 0
	}
	v, err := strconv.ParseUint(sValue, 10, 64)
	if err != nil {
		return 0
	}

	return v
}

func Uint64UserID(c *gin.Context) uint64 {
	return Uint64HeaderValue(c, HeaderUserID)
}

func Uint64TenantID(c *gin.Context) uint64 {
	return Uint64HeaderValue(c, HeaderTenantID)
}

package ginext

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	. "gitlab.com/goxp/cloud0/common"
)

// AuthRequiredMiddleware is required the request has to have x-user-id in header
// (it usually set by API Gateway)
func AuthRequiredMiddleware(c *gin.Context) {
	headers := struct {
		UserID   string `header:"x-user-id" validate:"required,min=1"`
		UserMeta string `header:"x-user-meta"`
		TenantID uint   `header:"x-tenant-id"`
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

func UintHeaderValue(c *gin.Context, headerName string) uint {
	sValue := c.GetHeader(headerName)
	if sValue == "" {
		return 0
	}
	v, err := strconv.Atoi(sValue)
	if err != nil {
		return 0
	}

	return uint(v)
}

func UintUserID(c *gin.Context) uint {
	return UintHeaderValue(c, HeaderUserID)
}

func UintTenantID(c *gin.Context) uint {
	return UintHeaderValue(c, HeaderTenantID)
}

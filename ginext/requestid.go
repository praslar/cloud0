package ginext

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "gitlab.com/goxp/cloud0/common"
)

func RequestIDMiddleware(c *gin.Context) {
	requestid := c.GetHeader(HeaderXRequestID)
	if requestid == "" {
		requestid = uuid.New().String()
		c.Header(HeaderXRequestID, requestid)
	}
	// set to context
	c.Set(HeaderXRequestID, requestid)

	c.Next()
}

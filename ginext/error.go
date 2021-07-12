package ginext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"gitlab.com/goxp/cloud0/logger"
)

var (
	_ ApiError = &apierr{}
)

type ApiError interface {
	Code() int
	MarshalJSON() ([]byte, error)
}

type apierr struct {
	code    int
	message string
}

func (e *apierr) Code() int {
	return e.code
}

func (e *apierr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"detail": e.message})
}

func (e *apierr) Error() string {
	return e.message
}

func NewError(code int, message string) error {
	return &apierr{code: code, message: message}
}

func ErrorHandler(c *gin.Context) {
	l := logger.WithCtx(c, "ErrorHandler")

	defer func() {
		if err := recover(); err != nil {
			l.WithField("debug.Stack", string(debug.Stack())).Warn("handle panic")

			switch v := err.(type) {
			case error:
				// unable to read the request body
				if v == io.EOF {
					_ = c.Error(NewError(http.StatusBadRequest, "invalid payload"))
					break
				}
				_ = c.Error(v)

			case string:
				_ = c.Error(NewError(http.StatusInternalServerError, v))
			default:
				_ = c.Error(NewError(http.StatusInternalServerError, fmt.Sprintf("unknown error: %v", v)))
			}
		}

		// no error
		if len(c.Errors) == 0 {
			return
		}

		l.WithField("errors.len", len(c.Errors)).Debug("handle stacked errors")

		for _, err := range c.Errors {
			l.WithError(err.Err).Debug("process error")
		}

		// just respond last error now
		err := c.Errors.Last().Err
		code := http.StatusInternalServerError
		if v, ok := err.(ApiError); ok {
			code = v.Code()
		}

		c.JSON(code, gin.H{"error": err})
	}()

	c.Next()

}

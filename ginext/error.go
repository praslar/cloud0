package ginext

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	var err error

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = NewError(http.StatusInternalServerError, fmt.Sprintf("unexpected error: %v", v))
			}
		}

		// no error
		if err == nil && len(c.Errors) == 0 {
			return
		}

		if len(c.Errors) > 0 {
			l = l.WithField("errors", c.Errors)
			if err != nil {
				l = l.WithField("recoveredError", err)
			} else {
				err = c.Errors.Last().Err
			}
		}

		l.Debug("handle request error")

		code := http.StatusInternalServerError
		if v, ok := err.(ApiError); ok {
			code = v.Code()
		}

		c.JSON(code, &GeneralBody{Error: err})
	}()

	c.Next()
}

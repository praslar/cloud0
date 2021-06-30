package ginext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/goxp/cloud0/log"
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
	l := log.Tag("ginext.ErrorHandler")
	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case error:
				// unable to read the request body
				if v == io.EOF {
					_ = c.Error(NewError(http.StatusBadRequest, "invalid payload"))
					break
				}
				_ = c.Error(v)

			case string:
				l.Debugf("erorr: %s", v)
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

		// we might need to log/write errors later, just respond last error now
		err := c.Errors.Last().Err

		l.WithError(err).Debug("handle error")

		code := http.StatusInternalServerError
		if v, ok := err.(ApiError); ok {
			code = v.Code()
		}

		c.JSON(code, gin.H{"error": err})
	}()

	c.Next()

}

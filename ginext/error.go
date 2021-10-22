package ginext

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"gitlab.com/praslar/cloud0/logger"
)

var (
	_ ApiError = &apiErr{}
)

// ApiError is an interface that supports return Code & Marshal to json
type ApiError interface {
	Code() int
	MarshalJSON() ([]byte, error)
}

type apiErr struct {
	code    int
	message string
}

func (e *apiErr) Code() int {
	return e.code
}

func (e *apiErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"detail": e.Error()})
}

func (e *apiErr) Error() string {
	return e.message
}

func NewError(code int, message string) error {
	return &apiErr{code: code, message: message}
}

func CreateErrorHandler(printStacks ...bool) gin.HandlerFunc {
	printStack := false
	if len(printStacks) > 0 {
		printStack = printStacks[0]
	}
	
	return func(c *gin.Context) {
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

			if err == nil && len(c.Errors) > 0 {
				err = c.Errors.Last().Err
			}

			l.WithError(err).Debug("handle request error")
			if printStack {
				fmt.Println(errors.Wrap(err, 1).ErrorStack())
			}

			code := http.StatusInternalServerError
			if v, ok := err.(ApiError); ok {
				code = v.Code()
			} else if v, ok := err.(*json.UnmarshalTypeError); ok {
				code = http.StatusBadRequest
				err = &validationErrors{
					fieldErrors: []ValidatorFieldError{
						&validatorFieldError{
							field:   v.Field,
							message: fmt.Sprintf("invalid type `%s`, requires `%s`", v.Value, v.Type.String()),
						},
					},
				}
			}

			c.JSON(code, &GeneralBody{Error: err})
		}()

		c.Next()
	}
}

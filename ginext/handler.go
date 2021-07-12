package ginext

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/goxp/cloud0/common"
)

type Request struct {
	GinCtx *gin.Context
	ctx    context.Context
}

type Response struct {
	Code   int
	Header http.Header
	*GeneralResponse
}

func NewResponse(code int) *Response {
	return &Response{
		Code:            code,
		GeneralResponse: &GeneralResponse{},
	}
}

func NewResponseData(code int, data interface{}) *Response {
	return &Response{
		Code:            code,
		GeneralResponse: NewPaginatedResponse(data, nil),
	}
}

func NewResponseWithPager(code int, data interface{}, pager *Pager) *Response {
	return &Response{
		Code:            code,
		GeneralResponse: NewPaginatedResponse(data, pager),
	}
}

type Handler func(r *Request) (*Response, error)

// NewRequest creates a new handler request
func NewRequest(c *gin.Context) *Request {
	ctx := FromGinRequestContext(c)
	req := &Request{
		GinCtx: c,
		ctx:    ctx,
	}

	return req
}

func (r *Request) Context() context.Context {
	return r.ctx
}

func WrapHandler(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			resp *Response
		)

		defer func() {
			if err != nil {
				_ = c.Error(err)
				return
			}

			// if any panic
			panicErr := recover()
			if panicErr != nil {
				_ = c.Error(err)
				return
			}

			for k, v := range resp.Header {
				for _, v_ := range v {
					c.Header(k, v_)
				}
			}

			if resp.Data != nil || resp.Error != nil {
				c.JSON(resp.Code, resp.GeneralResponse)
			} else {
				c.Status(resp.Code)
			}
		}()

		req := NewRequest(c)
		resp, err = handler(req)
	}
}

// MustBind does a binding on v with income request data
// it'll panic if any invalid data (and must be recovered by WrapHandler in this scope)
func (r *Request) MustBind(v interface{}) {
	err := r.GinCtx.ShouldBind(v)
	if err != nil {
		panic(err)
	}
}

func (r *Request) MustBindUri(v interface{}) {
	err := r.GinCtx.ShouldBindUri(v)
	if err != nil {
		panic(err)
	}
}

// MustNoError makes a ASSERT on err variable, panic when it's not nil
// then it must be recovered by WrapHandler
func (r *Request) MustNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func (r *Request) UintUserID() uint {
	return UintHeaderValue(r.GinCtx, common.HeaderUserID)
}

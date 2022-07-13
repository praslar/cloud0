package ginext

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/praslar/cloud0/common"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequiredMiddleware(t *testing.T) {

	cases := []struct {
		name       string
		userHeader string
		wantStatus int
	}{
		{
			name:       "MissingUserHeader",
			userHeader: "",
			wantStatus: 401,
		},
		{
			name:       "Success",
			userHeader: "10",
			wantStatus: 200,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("POST", "/protected", nil)
			c.Request.Header.Set(common.HeaderUserID, tc.userHeader)
			AuthRequiredMiddleware(c)
			c.Writer.WriteHeaderNow()

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetHeader(t *testing.T) {
	cases := []struct {
		name           string
		setHeaderName  string
		setHeaderValue string
		getHeaderName  string
		wantValue      interface{}
	}{
		{
			name:           "GetUintUserID",
			setHeaderName:  common.HeaderUserID,
			setHeaderValue: "10",
			getHeaderName:  common.HeaderUserID,
			wantValue:      uint64(10),
		},
		{
			name:           "GetUintTenantID",
			setHeaderName:  common.HeaderTenantID,
			setHeaderValue: "01",
			getHeaderName:  common.HeaderTenantID,
			wantValue:      uint64(1),
		},
		{
			name:          "Return0WhenMissing",
			getHeaderName: common.HeaderTenantID,
			wantValue:     uint64(0),
		},
		{
			name:           "Return0WhenInvalid",
			setHeaderName:  common.HeaderUserID,
			setHeaderValue: "abc",
			getHeaderName:  common.HeaderUserID,
			wantValue:      uint64(0),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("POST", "/protected", nil)
			if tc.setHeaderName != "" {
				c.Request.Header.Set(tc.setHeaderName, tc.setHeaderValue)
			}

			got := Uint64HeaderValue(c, tc.getHeaderName)
			assert.Equal(t, tc.wantValue, got)
		})
	}
}

func TestGetShortcutFunc(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/protected", nil)
	c.Request.Header.Set(common.HeaderUserID, "10")
	c.Request.Header.Set(common.HeaderTenantID, "01")
	c.Set(common.HeaderUserID, "10")

	assert.Equal(t, uint64(10), Uint64UserID(c))
	assert.Equal(t, "10", GetUserID(c))
	assert.Equal(t, uint64(1), Uint64TenantID(c))
}

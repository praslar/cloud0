package service

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"gitlab.com/goxp/cloud0/common"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	// setup
	_ = os.Setenv("PORT", "0")
	gin.SetMode(gin.TestMode)

	app := NewApp("echo", "v1")

	err := app.Initialize()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go app.Start(ctx)

	<-time.After(time.Millisecond * 100)

	req, err := http.NewRequest(http.MethodGet, "http://"+app.Listener().Addr().String() + "/status", nil)
	require.NoError(t, err)
	rsp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)

	assert.NotEmpty(t, rsp.Header.Get(common.HeaderXRequestID))
}

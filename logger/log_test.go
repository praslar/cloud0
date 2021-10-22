package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/praslar/cloud0/common"
)

type mockGetStringer func(string) string

func (f mockGetStringer) GetString(key string) string {
	return f(key)
}

func TestTagWithCtx(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "json")
	_ = os.Setenv("LOG_LEVEL", "debug")
	Init("test")
	entry := TagWithGetString("test", mockGetStringer(func(_ string) string {
		return "test-request-id"
	}))

	assert.Equal(t, "test", entry.Data["tag"])
	assert.Equal(t, "test-request-id", entry.Data[common.HeaderXRequestID])

	entry.Debug("finish log unit tests")
}

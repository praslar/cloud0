package ginext

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/praslar/cloud0/logger"
)

func TestMain(m *testing.M) {
	logger.Init("ginext.test")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

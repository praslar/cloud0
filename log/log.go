package log

import (
	"os"

	"github.com/sirupsen/logrus"
	. "gitlab.com/goxp/cloud0/common"
)

// GetStringer describe an object that has capacity to return a string via GetString
// It uses as Gin Context in case we want to cut off gin dependency here
type GetStringer interface {
	GetString(key string) string
}

func Init() {
	if l, e := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); e == nil {
		logrus.SetLevel(l)
	}
	if os.Getenv("LOG_FORMAT") == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "",
			DisableTimestamp: false,
			DataKey:          "",
			FieldMap:         nil,
			CallerPrettyfier: nil,
			PrettyPrint:      false,
		})
	}
	logrus.SetOutput(os.Stdout)
}

// Tag sets a tag name then returns a log entry ready to write
func Tag(tag string) *logrus.Entry {
	return logrus.WithField("tag", tag)
}

// TagWithCtx return a log entry from tag name & x-request-id in Gin context if has
func TagWithCtx(tag string, ctx GetStringer) *logrus.Entry {
	l := Tag(tag)
	if requestID := ctx.GetString(HeaderXRequestID); requestID != "" {
		l = l.WithField(HeaderXRequestID, requestID)
	}
	return l
}

package log

import (
	"context"
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

// TagWithGetString return a log entry from tag name & x-request-id in Gin context if has
func TagWithGetString(tag string, ctx GetStringer) *logrus.Entry {
	l := Tag(tag)
	if requestID := ctx.GetString(HeaderXRequestID); requestID != "" {
		l = l.WithField(HeaderXRequestID, requestID)
	}
	return l
}

func WithCtx(ctx context.Context, tag string) *logrus.Entry {
	l := Tag(tag)
	if requestID := ctx.Value("x-request-id").(string); requestID != "" {
		l = l.WithField("x-request-id", requestID)
	}
	return l
}

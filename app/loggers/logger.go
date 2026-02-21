package loggers

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var Log *logrus.Logger

func FromContext(ctx context.Context) *logrus.Entry {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return Log.WithField("trace_id", span.SpanContext().TraceID().String())
	}
	return Log.WithFields(logrus.Fields{})
}

func InitLogger(env string) error {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	// Set log level from environment variable if provided
	levelStr := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(levelStr)
	if err == nil {
		Log.SetLevel(level)
	} else {
		if env == "production" {
			Log.SetLevel(logrus.InfoLevel)
		} else {
			Log.SetLevel(logrus.DebugLevel)
		}
	}

	if env == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return nil
}

func Sync() {
	// Logrus doesn't require explicit syncing for stdout
}

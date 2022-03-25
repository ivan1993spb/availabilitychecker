package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey struct{}

var loggerKey = contextKey{}

func NewContext(parent context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(parent, loggerKey, logger)
}

func FromContext(ctx context.Context) *logrus.Entry {
	if logger, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return logger
	}

	return DiscardEntry
}

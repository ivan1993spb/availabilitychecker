package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Entry {
	logger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	return logrus.NewEntry(logger)
}

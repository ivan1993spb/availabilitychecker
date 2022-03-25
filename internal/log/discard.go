package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var DiscardLogger = &logrus.Logger{
	Out:          io.Discard,
	Hooks:        make(logrus.LevelHooks),
	Formatter:    new(NullFormatter),
	Level:        logrus.PanicLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

var DiscardEntry = logrus.NewEntry(DiscardLogger)

type NullFormatter struct{}

func (*NullFormatter) Format(*logrus.Entry) ([]byte, error) {
	return nil, nil
}

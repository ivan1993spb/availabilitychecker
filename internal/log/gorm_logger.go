package log

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SkipErrRecordNotFound: true,
	}
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	FromContext(ctx).WithContext(ctx).WithField("module", "gorm").Infof(s, args)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	FromContext(ctx).WithContext(ctx).WithField("module", "gorm").Warnf(s, args)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	FromContext(ctx).WithContext(ctx).WithField("module", "gorm").Errorf(s, args)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time,
	fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	sql, _ := fc()

	entry := FromContext(ctx)
	fields := logrus.Fields{
		"module":  "gorm",
		"elapsed": elapsed,
	}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) &&
		l.SkipErrRecordNotFound) {

		fields[logrus.ErrorKey] = err
		entry.WithContext(ctx).WithFields(fields).Error(sql)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		entry.WithContext(ctx).WithFields(fields).Warn(sql)
		return
	}

	entry.WithContext(ctx).WithFields(fields).Debug(sql)
}

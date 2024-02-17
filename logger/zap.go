package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.SugaredLogger
	once     sync.Once
)

// Initialize logger instance.
func InitLogger(isDev bool) error {
	var err error
	once.Do(func() {
		var logger *zap.Logger

		if isDev {
			logger, err = zap.NewDevelopment(
				zap.AddCaller(),
				zap.AddCallerSkip(1),
			)
		} else {
			logger, err = zap.NewProduction(
				zap.AddCaller(),
				zap.AddCallerSkip(1),
			)
		}
		instance = logger.Sugar()
	})
	return err
}

// / ---
// / Wrapper Function
// / ----
func GetLogger() *zap.SugaredLogger {
	once.Do(func() {
		if instance == nil {
			InitLogger(false)
		}
	})
	return instance
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	instance.Errorf(template, args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Infof(template string, args ...interface{}) {
	instance.Infof(template, args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	instance.Warnf(template, args...)
}

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	instance.Debugf(template, args...)
}

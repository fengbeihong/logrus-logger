package log

import (
	"time"

	"strings"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	UTCTime   = "utc"
	LocalTime = "local"
)

var DefaultMyLogger *MyLogger

func InitDefaultMyLogger(level logrus.Level, path string) {
	DefaultMyLogger = NewMyLogger(level, path)
}

func SetSkip(skip int) {
	DefaultMyLogger.SetSkip(skip)
}

func SetLogLevel(level logrus.Level) {
	DefaultMyLogger.logger.Level = level
}

func SetLogFormatter(formatter logrus.Formatter) {
	DefaultMyLogger.logger.Formatter = formatter
}

func SetHooks(path string) {
	DefaultMyLogger.SetHooks(path)
}

func SetLogMaxAge(maxAge time.Duration) {
	DefaultMyLogger.maxAge = maxAge
}

func SetRotationTime(rotationTime time.Duration) {
	DefaultMyLogger.rotationTime = rotationTime
}

func SetClockTime(clockTime string) {
	DefaultMyLogger.SetClockTime(clockTime)
}

func Debug(args ...interface{}) {
	DefaultMyLogger.Debug(args...)
}

func Info(args ...interface{}) {
	DefaultMyLogger.Info(args...)
}

func Warn(args ...interface{}) {
	DefaultMyLogger.Warn(args...)
}

func Error(args ...interface{}) {
	DefaultMyLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	DefaultMyLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	DefaultMyLogger.Panic(args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultMyLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultMyLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultMyLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultMyLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultMyLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	DefaultMyLogger.Panicf(format, args...)
}

type MyLogger struct {
	logger       *logrus.Logger
	skip         int
	logDir       string
	logFile      string
	clockTime    rotatelogs.Clock
	maxAge       time.Duration
	rotationTime time.Duration
}

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

// NewMyLogger Return Default MyLogger struct
func NewMyLogger(level logrus.Level, path string) *MyLogger {
	var ml = &MyLogger{}
	ml.logger = logrus.New()

	ml.SetLogLevel(level)
	ml.SetLogFormatter(&Formatter{})
	ml.SetClockTime(LocalTime)
	ml.SetLogMaxAge(time.Duration(86400) * time.Second)
	ml.SetRotationTime(time.Duration(604800) * time.Second)

	if strings.HasSuffix(path, ".log") {
		path = path[:strings.LastIndex(path, ".log")]
	}
	ml.SetHooks(path)

	return ml
}

func (ml *MyLogger) SetSkip(skip int) {
	ml.skip = skip
}

func (ml *MyLogger) SetLogLevel(level logrus.Level) {
	ml.logger.Level = level
}

func (ml *MyLogger) SetLogFormatter(formatter logrus.Formatter) {
	ml.logger.Formatter = formatter
}

func (ml *MyLogger) SetLogMaxAge(maxAge time.Duration) {
	ml.maxAge = maxAge
}

func (ml *MyLogger) SetRotationTime(rotationTime time.Duration) {
	ml.rotationTime = rotationTime
}

func (ml *MyLogger) SetClockTime(clockTime string) {
	switch clockTime {
	case UTCTime:
		ml.clockTime = rotatelogs.UTC
	case LocalTime:
		ml.clockTime = rotatelogs.Local
	}
}

func (ml *MyLogger) SetHooks(path string) {
	errorWriter, _ := rotatelogs.New(
		path+".error.%Y-%m-%d.log",
		rotatelogs.WithLinkName(path+".error.log"),
		rotatelogs.WithClock(ml.clockTime),
		rotatelogs.WithMaxAge(ml.maxAge),
		rotatelogs.WithRotationTime(ml.rotationTime),
	)
	infoWriter, _ := rotatelogs.New(
		path+".info.%Y-%m-%d.log",
		rotatelogs.WithLinkName(path+".info.log"),
		rotatelogs.WithClock(ml.clockTime),
		rotatelogs.WithMaxAge(ml.maxAge),
		rotatelogs.WithRotationTime(ml.rotationTime),
	)

	debugWriter, _ := rotatelogs.New(
		path+".debug.%Y-%m-%d.log",
		rotatelogs.WithLinkName(path+".debug.log"),
		rotatelogs.WithClock(ml.clockTime),
		rotatelogs.WithMaxAge(ml.maxAge),
		rotatelogs.WithRotationTime(ml.rotationTime),
	)

	ml.logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.PanicLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.WarnLevel:  infoWriter,
			logrus.InfoLevel:  infoWriter,
			logrus.DebugLevel: debugWriter,
		},
		ml.logger.Formatter,
	))
}

// Debug logs a message at level Debug on the standard ml.logger.
func (ml *MyLogger) Debug(args ...interface{}) {
	if ml.logger.Level >= logrus.DebugLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Debug(args...)
	}
}

// Info logs a message at level Info on the standard ml.logger.
func (ml *MyLogger) Info(args ...interface{}) {
	if ml.logger.Level >= logrus.InfoLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Info(args...)
	}
}

// Warn logs a message at level Warn on the standard ml.logger.
func (ml *MyLogger) Warn(args ...interface{}) {
	if ml.logger.Level >= logrus.WarnLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Warn(args...)
	}
}

// Error logs a message at level Error on the standard ml.logger.
func (ml *MyLogger) Error(args ...interface{}) {
	if ml.logger.Level >= logrus.ErrorLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Error(args...)
	}
}

// Fatal logs a message at level Fatal on the standard ml.logger.
func (ml *MyLogger) Fatal(args ...interface{}) {
	if ml.logger.Level >= logrus.FatalLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Fatal(args...)
	}
}

// Panic logs a message at level Panic on the standard ml.logger.
func (ml *MyLogger) Panic(args ...interface{}) {
	if ml.logger.Level >= logrus.PanicLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Panic(args...)
	}
}

// Debug logs a message at level Debug on the standard ml.logger.
func (ml *MyLogger) Debugf(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.DebugLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Debugf(format, args...)
	}
}

// Info logs a message at level Info on the standard ml.logger.
func (ml *MyLogger) Infof(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.InfoLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Infof(format, args...)
	}
}

// Warn logs a message at level Warn on the standard ml.logger.
func (ml *MyLogger) Warnf(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.WarnLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Warnf(format, args...)
	}
}

// Error logs a message at level Error on the standard ml.logger.
func (ml *MyLogger) Errorf(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.ErrorLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Errorf(format, args...)
	}
}

// Fatal logs a message at level Fatal on the standard ml.logger.
func (ml *MyLogger) Fatalf(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.FatalLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Fatalf(format, args...)
	}
}

// Panic logs a message at level Panic on the standard ml.logger.
func (ml *MyLogger) Panicf(format string, args ...interface{}) {
	if ml.logger.Level >= logrus.PanicLevel {
		entry := ml.logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Panicf(format, args...)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) InfoWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.InfoLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Info(l)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) DebugWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.DebugLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Debug(l)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) WarnWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.WarnLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Warn(l)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) ErrorWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.ErrorLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Error(l)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) FatalWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.FatalLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Fatal(l)
	}
}

// Debug logs a message with fields at level Debug on the standard ml.logger.
func (ml *MyLogger) PanicWithFields(l interface{}, f Fields) {
	if ml.logger.Level >= logrus.PanicLevel {
		entry := ml.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(ml.skip)
		entry.Panic(l)
	}
}

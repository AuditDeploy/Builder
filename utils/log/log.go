package log

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	defaultLogLevel := zapcore.Level(logLevel)

	layout := "01-02-2006"

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	globalPath, _ := os.LookupEnv("GLOBAL_LOGS_PATH")

	t := time.Now()
	logfile, _ := os.OpenFile(filepath.Join(globalPath, t.Format(layout), ".json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logfile)
	fileEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zap.ReplaceGlobals(logger)
}

func Debug(msg string, args ...interface{}) {
	if logLevel <= DEBUG {
		if len(args) > 0 {
			zap.S().Debugf(msg, args...)
		} else {
			zap.S().Debug(msg)
		}
	}
}

func Info(msg string, args ...interface{}) {
	if logLevel <= INFO {
		if len(args) > 0 {
			zap.S().Infof(msg, args...)
		} else {
			zap.S().Info(msg)
		}
	}
}

func Warn(msg string, args ...interface{}) {
	if logLevel <= WARNING {
		if len(args) > 0 {
			zap.S().Warnf(msg, args...)
		} else {
			zap.S().Warn(msg)
		}
	}
}

func Error(msg string, args ...interface{}) {
	if logLevel <= ERROR {
		if len(args) > 0 {
			zap.S().Errorf(msg, args...)
		} else {
			zap.S().Error(msg)
		}
	}
}

func Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		zap.S().Fatalf(msg, args...)
	} else {
		zap.S().Fatal(msg)
	}
}

func SetLevel(level Level) {
	level = level
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	ErrorFATAL
)

var logLevel Level

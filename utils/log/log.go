package log

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func NewLogger(logFileName string, path string) *zap.Logger {
        defaultLogLevel := zapcore.Level(logLevel)

        config := zap.NewProductionEncoderConfig()
        config.TimeKey = "timestamp"
        config.EncodeTime = zapcore.ISO8601TimeEncoder

        logfile, _ := os.OpenFile(filepath.Join(path, logFileName + ".json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

        writer := zapcore.AddSync(logfile)
        fileEncoder := zapcore.NewJSONEncoder(config)
        core := zapcore.NewTee(
                zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
        )

        logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    return logger
}

func init() {
        args := os.Args[1:]

	// If verbose flag given display Builder logs to console
	for i := 0; i < len(args); i++ {
		if args[i] == "-v" || args[i] == "--verbose" {
        		logger, _ := zap.NewDevelopment()

        		zap.ReplaceGlobals(logger)
		}
	}
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

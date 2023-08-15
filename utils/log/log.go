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

	logfile, _ := os.OpenFile(filepath.Join(path, logFileName+".json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logfile)
	fileEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	// If debug flag given display caller in log as well
	args := os.Args[1:]
	debugFlag := false
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" || args[i] == "--debug" {
			debugFlag = true
		}
	}

	if debugFlag {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger
}

func init() {
	args := os.Args[1:]

	// If verbose flag given display Builder logs to console
	for i := 0; i < len(args); i++ {
		if args[i] == "-v" || args[i] == "--verbose" {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			zap.ReplaceGlobals(logger)
		}
	}
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

package log

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func NewLogger(logFileName string, path string) (*zap.Logger, func()) {
	defaultLogLevel := zapcore.Level(logLevel)

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	//logfile, _ := os.OpenFile(filepath.Join(path, logFileName+".json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer, closeFile, err := zap.Open(filepath.Join(path, logFileName+".json"))
	if err != nil {
		fmt.Println("logger err")
	}

	// If debug flag given display caller in log and print build logs to console
	args := os.Args[1:]
	verboseFlag := false
	debugFlag := false
	for i := 0; i < len(args); i++ {
		if args[i] == "-v" || args[i] == "--verbose" {
			verboseFlag = true
		}
		if args[i] == "-d" || args[i] == "--debug" {
			debugFlag = true
		}
	}

	var core zapcore.Core

	// If verbose flag given display build logs to console as well as to file
	if verboseFlag {
		fileEncoder := zapcore.NewJSONEncoder(config)
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	} else {
		fileEncoder := zapcore.NewJSONEncoder(config)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		)
	}

	// If debug flag given add Builder caller to build logs
	if debugFlag {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger, closeFile
}

func init() {
	args := os.Args[1:]

	// If debug flag given display Builder logs to console
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" || args[i] == "--debug" {
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

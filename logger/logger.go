package logger

import (
	"log"
	"os"
	"strings"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func CreateLogs(filePath string) {
	logsDir := os.Getenv("BUILDER_LOGS_DIR")

	if os.Getenv("BUILDER_COMMAND") == "true" {
		// _, extName := artifact.ExtExistsFunction(logsDir, ".txt")
		file, err := os.OpenFile("hi", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		
		InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	} else {
	//log file name = parentDir
	newFileName := filePath[strings.LastIndex(filePath, "/")+1:]
	file, err := os.OpenFile(logsDir+"/"+newFileName+"_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

}

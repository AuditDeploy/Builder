package directory

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
)

func MakeLogDir(parentDirPath string) {
	logDirPath := parentDirPath + "/logs"

	logDir(logDirPath)
	logger.CreateLogs(logDirPath)
}

func logDir(logDirPath string) {
	//check if logDir path exists
	_, err := os.Stat(logDirPath)

	if err == nil {
		fmt.Println("Path already exists")
		logger.WarningLogger.Println("Path already exists")
	}

	if os.IsNotExist(err) {
		errDir := os.Mkdir(logDirPath, 0755)

		if errDir != nil {
			log.Fatal(err)
		}

	}

	//check logsDir env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_LOGS_DIR")
	if !present {
		os.Setenv("BUILDER_LOGS_DIR", logDirPath)
	} else {
		fmt.Println("BUILDER_LOGS_DIR", val)
	}
}

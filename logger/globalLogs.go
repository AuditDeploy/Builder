package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

//creates global logs
func GlobalLogs() {
	var globalLogsPath string
	//look for global logs env var, add
	val, present := os.LookupEnv("GLOBAL_LOGS_PATH")
	if !present {
		//create global logs if it does not exists
		if _, err := os.Stat("./globalLogs"); os.IsNotExist(err) {
			os.Mkdir("./globalLogs", 0755)
		}
		globalLogsPath = "./globalLogs/logs.txt"
	} else {
		globalLogsPath = val
	}

	file, err := os.OpenFile(globalLogsPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	//copy log over to global logs
	logsDir := os.Getenv("BUILDER_LOGS_DIR")
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	fileName := parentDir[strings.LastIndex(parentDir, "/")+1:]

	//append log to global logs
	from, err := os.Open(logsDir + "/" + fileName + "_logs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(globalLogsPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

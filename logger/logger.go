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

	//points back to already created log.txt if using 'builder' cmd
	// if os.Getenv("BUILDER_COMMAND") == "true" {
	// 	path, _ := os.Getwd()
	// 	fmt.Println(path)

	// 	var dirPath string
	// 	if strings.Contains(path, "workspace") && strings.Contains(path, "temp") {
	// 		dirPath = strings.Replace(path, "/workspace/temp", "", 1)
	// 		fmt.Println("IN TEMP: ", dirPath)

	// 	} else if strings.Contains(path, "workspace") {
	// 		dirPath = strings.TrimRight(path, "/workspace")
	// 		fmt.Println("IN WRKSPC: ", dirPath)
	// 	}

	// 	_, extName := artifact.ExtExistsFunction(dirPath+"/logs/", ".txt")
	// 	file, err := os.OpenFile(dirPath+"/logs/"+extName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// 	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// 	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// } else {
	//log file name = parentDir, creates new logs.txt
	newFileName := filePath[strings.LastIndex(filePath, "/")+1:]
	file, err := os.OpenFile(logsDir+"/"+newFileName+"_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	// }

}

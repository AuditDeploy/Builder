package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//creates global logs
func GlobalLogs() {
	//create global logs if it does not exists
	if _, err := os.Stat("./globalLogs"); os.IsNotExist(err) {
		os.Mkdir("./globalLogs", 0755)

	}

	file, err := os.OpenFile("./globalLogs/logs.txt", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	//create global logs environment
	val, present := os.LookupEnv("BUILDER_GLOBAL_LOGS_DIR")
	if !present {
		os.Setenv("BUILDER_GLOBAL_LOGS_DIR", "./globalLogs")
	} else {
		fmt.Println("BUILDER_GLOBAL_LOGS_DIR", val)
	}

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

	to, err := os.OpenFile("./globalLogs/logs.txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

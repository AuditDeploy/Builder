package logger

import (
	"fmt"
	"os"
	"os/exec"
)

//creates global logs
func GlobalLogs() {
	//create global logs if it does not exists
	if _, err := os.Stat("./globalLogs"); os.IsNotExist(err) {
		os.Mkdir("./globalLogs", 0755)
	}

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

	fileName := ConvertParentPathToFileName(parentDir)

	exec.Command("cp", logsDir+"/"+fileName+".txt", "./globalLogs").Run()
}

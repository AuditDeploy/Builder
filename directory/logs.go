package directory

import (
	"Builder/spinner"
	"fmt"
	"os"
)

func logDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		spinner.LogMessage("Path already exists", "info")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {

		errDir := os.Mkdir(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			spinner.LogMessage("failed to make directory: "+err.Error(), "fatal")
		}

	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_LOGS_DIR")
	if !present {
		os.Setenv("BUILDER_LOGS_DIR", path)
	} else {
		fmt.Println("BUILDER_LOGS_DIR", val)
	}
	return true, err
}

// MakeLogsDir does...
func MakeLogsDir(path string) {
	visiblePath := path + "/logs"
	logDir(visiblePath)
}

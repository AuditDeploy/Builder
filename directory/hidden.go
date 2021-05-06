package directory

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
)

func MakeHiddenDir(parentDirPath string) {
	hiddenDirPath := parentDirPath + "/.hidden"
	hiddenDir(hiddenDirPath)
}

func hiddenDir(hiddenDirPath string) {
	//check if hiddenDir path exists
	_, err := os.Stat(hiddenDirPath)

	if err == nil {
		fmt.Println("Path already exists")
		logger.WarningLogger.Println("Path already exists")

	}

	if os.IsNotExist(err) {
		errDir := os.Mkdir(hiddenDirPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	//check hiddenDir env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_HIDDEN_DIR")
	if !present {
		os.Setenv("BUILDER_HIDDEN_DIR", hiddenDirPath)
	} else {
		fmt.Println("BUILDER_HIDDEN_DIR", val)
	}
}

package directory

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
)

func MakeWorkspaceDir(parentDirPath string) {

	workDirPath := parentDirPath + "/workspace"
	workSpaceDir(workDirPath)
}

func workSpaceDir(workDirPath string) {
	//check if file path exists
	_, err := os.Stat(workDirPath)

	if err == nil {
		fmt.Println("Path already exists")
		logger.WarningLogger.Println("Path already exists")
	}

	if os.IsNotExist(err) {
		errDir := os.Mkdir(workDirPath, 0755)

		if errDir != nil {
			log.Fatal(err)
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_WORKSPACE_DIR")
	if !present {
		os.Setenv("BUILDER_WORKSPACE_DIR", workDirPath)
	} else {
		fmt.Println("BUILDER_WORKSPACE_DIR", val)
	}

}

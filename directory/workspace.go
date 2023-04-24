package directory

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"runtime"
)

func workSpaceDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		logger.WarningLogger.Println("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {

		errDir := os.Mkdir(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			log.Fatal(err)
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_WORKSPACE_DIR")
	if !present {
		os.Setenv("BUILDER_WORKSPACE_DIR", path)
	} else {
		fmt.Println("BUILDER_WORKSPACE_DIR", val)
	}
	return true, err
}

// MakeWorkspaceDir does...
func MakeWorkspaceDir(path string) {

	//workPath := path + "/workspace"
	//testing create workspace dir
	var pathSeperator = "/"

	if runtime.GOOS == "windows" {
		pathSeperator = "\\"
	}
	/* currentDirectory, _ := os.Getwd()
	logger.InfoLogger.Println(filepath.Join(currentDirectory, path, "workspace"))
	workPath := filepath.Join(currentDirectory, path, "workspace") */
	workPath := path + pathSeperator + "workspace"

	workSpaceDir(workPath)

}

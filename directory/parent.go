package directory

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"Builder/utils"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()

// MakeDirs does...
func MakeDirs() {
	//handles -n flag
	name := utils.GetName()

	//add Unix timestamp to dir name
	currentTime := time.Now().Unix()

	//check for projectPath env from builder.yaml
	configPath := os.Getenv("BUILDER_DIR_PATH")

	unixTime := strconv.FormatInt(currentTime, 10)
	os.Setenv("BUILDER_TIMESTAMP", unixTime)

	var path string
	if configPath != "" {
		// used for 'config' cmd, set by builder.yaml
		path = configPath + "/" + name + "_" + unixTime
	} else {
		// local path, used for 'init' cmd/default
		path = "./" + name + "_" + unixTime
	}

	MakeParentDir(path)

	MakeHiddenDir(path)
	MakeWorkspaceDir(path)

	MakeLogsDir(path)
	MakeBuilderDir()
}

func MakeParentDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		BuilderLog.Info("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			BuilderLog.Fatalf("failed to create directory", path, err)
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_PARENT_DIR")
	if !present {
		os.Setenv("BUILDER_PARENT_DIR", path)
	} else {
		fmt.Println("BUILDER_PARENT_DIR", val)
	}

	return true, err
}

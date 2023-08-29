package directory

import (
	"fmt"
	"os"
	"strings"
	"time"

	"Builder/spinner"
	"Builder/utils"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()

// MakeDirs does...
func MakeDirs() {
	//handles -n flag
	name := utils.GetName()

	//check for projectPath env from builder.yaml
	configPath := os.Getenv("BUILDER_DIR_PATH")

	var path string
	if configPath != "" {
		// used for 'config' cmd, set by builder.yaml
		path = configPath + "/" + name + "_" + name
	} else {
		// local path, used for 'init' cmd/default
		path = "./" + name + "_" + name
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
		spinner.LogMessage("Path already exists", "info")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			spinner.LogMessage("failed to create directory at "+path+": "+err.Error(), "fatal")
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

func UpdateParentDirName(pathWithWrongParentName string) string {
	oldName, _ := os.LookupEnv("BUILDER_PARENT_DIR")
	projectName := utils.GetName()
	startTime, _ := time.Parse(time.RFC850, os.Getenv("BUILD_START_TIME"))
	unixTimestamp := startTime.Unix()
	newName := strings.TrimSuffix(oldName, projectName) + fmt.Sprint(unixTimestamp)

	path := os.Getenv("BUILDER_DIR_PATH")

	if oldName[0:2] == "./" {
		err := os.Rename(path+"/"+oldName[2:], path+"/"+newName[2:])
		if err != nil {
			fmt.Println(err.(*os.LinkError).Err)
			spinner.LogMessage("could not rename parent dir", "fatal")
		}
	} else {
		err := os.Rename(oldName, newName)
		if err != nil {
			fmt.Println(err.(*os.LinkError).Err)
			spinner.LogMessage("could not rename parent dir", "fatal")
		}
	}

	os.Setenv("BUILDER_PARENT_DIR", newName)
	os.Setenv("BUILDER_WORKSPACE_DIR", newName+"/workspace")
	os.Setenv("BUILDER_LOGS_DIR", newName+"/logs")

	// Return new path with new parent directory name
	newPath := strings.Replace(pathWithWrongParentName, oldName[2:], newName[2:], 1)

	return newPath
}

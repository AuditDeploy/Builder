package directory

import (
	"fmt"
	"os"
	"strings"
	"time"

	"Builder/utils"
	"Builder/utils/log"

	"github.com/manifoldco/promptui"
)

// MakeDirs does...
func MakeDirs() {
	//handles -n flag
	name := utils.GetName()

	//check for projectPath env from builder.yaml
	configPath := os.Getenv("BUILDER_DIR_PATH")

	var path string
	if configPath != "" {
		// used for 'config' cmd, set by builder.yaml
		path = configPath + "/" + name + "_START"
	} else {
		// local path, used for 'init' cmd/default
		path = "./" + name + "_START"
	}

	MakeParentDir(path)

	MakeHiddenDir(path)
	MakeWorkspaceDir(path)
}

func MakeParentDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		log.Info("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		//bypass Prompt msg
		if bypassPrompt() {
			errDir := os.MkdirAll(path, 0755)
			//should return nil once directory is made, if not, throw err
			if errDir != nil {
				log.Fatal("failed to create directory", path, err)
			}
		} else {
			//prompt user if they'd like dir to be created
			mk := yesNo()

			if mk {
				errDir := os.MkdirAll(path, 0755)
				//should return nil once directory is made, if not, throw err
				if errDir != nil {
					log.Fatal("failed to create directory", path, err)
				}

			} else {
				//logger.ErrorLogger.Println("Please create a directory for the Builder")
				log.Fatal("Please create a directory for the Builder")
				return true, err
			}
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
	startTime, _ := time.Parse(time.RFC850, os.Getenv("BUILD_START_TIME"))
	unixTimestamp := startTime.Unix()
	newName := strings.TrimSuffix(oldName, "START") + fmt.Sprint(unixTimestamp)

	path := os.Getenv("BUILDER_DIR_PATH")

	err := os.Rename(path+"/"+oldName[2:], path+"/"+newName[2:])
	if err != nil {
		log.Fatal("could not rename parent dir")
	}

	os.Setenv("BUILDER_PARENT_DIR", newName)

	// Return new path with new parent directory name
	newPath := strings.Replace(pathWithWrongParentName, oldName[2:], newName[2:], 1)

	return newPath
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Create A Directory? [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func bypassPrompt() bool {
	args := os.Args[1:]

	yesFlag := false

	val := os.Getenv("BYPASS_PROMPTS")
	if val == "true" {
		yesFlag = true
	}
	for _, val := range args {
		if val == "--yes" || val == "-y" {
			yesFlag = true
		}
	}

	return yesFlag
}

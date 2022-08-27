package directory

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"builder/logger"
	"builder/utils"

	"github.com/manifoldco/promptui"
)

//MakeDirs does...
func MakeDirs() {
	//handles -n flag
	name := utils.GetName()

	//add Unix timestamp to dir name
	currentTime := time.Now().Unix()

	//check for projectPath env from builder.yaml
	configPath := os.Getenv("builder_DIR_PATH")

	unixTime := strconv.FormatInt(currentTime, 10)
	os.Setenv("builder_TIMESTAMP", unixTime)

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
	MakeLogDir(path)
	MakeWorkspaceDir(path)
}

func MakeParentDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		logger.WarningLogger.Println("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		//bypass Prompt msg
		if bypassPrompt() {
			errDir := os.MkdirAll(path, 0755)
			//should return nil once directory is made, if not, throw err
			if errDir != nil {
				log.Fatal(err)
			}
		} else {
			//prompt user if they'd like dir to be created
			mk := yesNo()

			if mk {
				errDir := os.MkdirAll(path, 0755)
				//should return nil once directory is made, if not, throw err
				if errDir != nil {
					log.Fatal(err)
				}

			} else {
				//logger.ErrorLogger.Println("Please create a directory for the builder")
				log.Fatal("Please create a directory for the builder")
				return true, err
			}
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("builder_PARENT_DIR")
	if !present {
		os.Setenv("builder_PARENT_DIR", path)
	} else {
		fmt.Println("builder_PARENT_DIR", val)
	}

	return true, err
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Create A Directory? [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func bypassPrompt() bool {
	args := os.Args[1:]

	yesFlag := false

	for _, val := range args {
		if val == "--yes" || val == "-y" {
			yesFlag = true
		}
	}

	return yesFlag
}

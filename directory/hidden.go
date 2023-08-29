package directory

import (
	"Builder/spinner"
	"Builder/utils"
	"fmt"
	"os"
	"strings"
)

func hiddenDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		spinner.LogMessage("Path already exists", "warn")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			spinner.LogMessage("failed to create hidden directory: "+err.Error(), "fatal")
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_HIDDEN_DIR")
	if !present {
		os.Setenv("BUILDER_HIDDEN_DIR", path)
	} else {
		fmt.Println("BUILDER_HIDDEN_DIR", val)
	}
	return true, err
}

// MakeHiddenDir does...
func MakeHiddenDir(path string) {

	if os.Getenv("HIDDEN_DIR_ENABLED") == "true" {
		hiddenPath := path + "/.hidden"
		hiddenDir(hiddenPath)
	} else {
		repo := utils.GetRepoURL()
		repoName := strings.TrimSuffix(repo[strings.LastIndex(repo, "/"):], ".git")
		// Don't add extra slash if one exists
		var visiblePath string
		if strings.Contains(repoName, "/") {
			visiblePath = path + repoName
		} else {
			visiblePath = path + "/" + repoName
		}

		hiddenDir(visiblePath)
	}
}

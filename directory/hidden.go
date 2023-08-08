package directory

import (
	"Builder/utils/log"
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
		log.Warn("Path already exists")

	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			log.Fatal("failed to create hidden directory", err)
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
		visiblePath := path + "/" + repoName
		hiddenDir(visiblePath)
	}
}

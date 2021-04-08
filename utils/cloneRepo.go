package utils

import (
	"Builder/logger"
	"fmt"
	"os"
	"os/exec"
)

//CloneRepo grabs url
func CloneRepo() {

	repo := GetRepoURL()
	
	
	//clone repo with url from args
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	//if config cmd clone to temp dir
	if (hiddenDir == "") {
		errDir := os.Mkdir("./tempRepo", 0755)
		if errDir != nil {
			fmt.Println(errDir)
		}	
		cmd := exec.Command("git", "clone", repo, "./tempRepo")
		cmd.Run()
	} else {
		//if init cmd, clone to hidden dir
		cmd := exec.Command("git", "clone", repo, hiddenDir)
		logger.InfoLogger.Println(cmd)
		cmd.Run()
	}
}

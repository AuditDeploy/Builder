package utils

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//CloneRepo grabs url
func CloneRepo() {

	repo := GetRepoURL()

	//clone repo with url from args
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	bFlagExists := BranchFlagExists()
	branchName := GetBranchName()
	allRepoBranches := GetAllRepoBranches()
	branchNameExists := BranchNameExists(allRepoBranches)

	//if config cmd clone to temp dir
	if hiddenDir == "" {
		errDir := os.Mkdir("./tempRepo", 0755)
		if errDir != nil {
			fmt.Println(errDir)
		}

		if bFlagExists {
			if branchNameExists {
				cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, "./tempRepo")
				cmd.Run()
			} else {
				log.Fatal("Branch does not exists")
			}

		} else {
			cmd := exec.Command("git", "clone", repo, "./tempRepo")
			cmd.Run()
		}

		cmd := exec.Command("git", "clone", repo, "./tempRepo")
		cmd.Run()
	} else {

		if bFlagExists {
			if branchNameExists {
				cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, hiddenDir)
				cmd.Run()
				logger.InfoLogger.Println(cmd)
			} else {
				log.Fatal("Branch does not exists")
			}

		} else {
			cmd := exec.Command("git", "clone", repo, hiddenDir)
			cmd.Run()
			logger.InfoLogger.Println(cmd)
		}

	}
}

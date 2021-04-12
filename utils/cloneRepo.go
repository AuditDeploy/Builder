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

	//if config cmd clone to temp dir
	if hiddenDir == "" {
		errDir := os.Mkdir("./tempRepo", 0755)
		if errDir != nil {
			fmt.Println(errDir)
		}
		cmd := exec.Command("git", "clone", repo, "./tempRepo")
		cmd.Run()
	} else {
		//if init cmd, clone to hidden dir
		bFlagExists, branchName := CloneBranch()

		//check if branch exist
		branches, _, _, _ := GitHashAndName()
		branchExists, _ := BranchNameExists(branches)

		if bFlagExists {
			if branchExists {
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

func CloneBranch() (bool, string) {
	args := os.Args[1:]

	//if branch is empty string strings.Contain does not work, function found in metadata
	branchName := "%$F"
	branchFlag := false

	for i, v := range args {
		if v == "-b" || v == "--branch" {
			if len(args) <= i+1 {
				logger.ErrorLogger.Println("No Repo Url Provided")
				log.Fatal("No Branch Name Provided")

			} else {
				branchName = args[i+1]
				branchFlag = true
			}
		}
	}

	return branchFlag, branchName
}

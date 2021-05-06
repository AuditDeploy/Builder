package utils

import (
	"Builder/logger"
	"log"
	"os"
	"os/exec"
	"strings"
)

func BranchFlagExists() bool {
	bFlagExists, _ := GetBranchFlagAndName()
	return bFlagExists
}

func GetBranchName() string {
	_, branchName := GetBranchFlagAndName()
	return branchName
}

func BranchNameExists(branches []string) bool {
	branchExists := false
	branchName := GetBranchName()

	for _, branch := range branches {
		if branch[strings.LastIndex(branch, "/")+1:] == branchName {
			branchExists = true
		}
	}

	return branchExists
}

func GetAllRepoBranches() []string {
	repo := GetRepoURL()
	output, _ := exec.Command("git", "ls-remote", repo).Output()
	stringifyAllBranches := string(output)
	allRepoBranches := strings.Split(stringifyAllBranches, "\n")

	return allRepoBranches
}

func GetBranchGitHash() string {
	allBranches := GetAllRepoBranches()
	branchExists := BranchNameExists(allBranches)

	branchName := GetBranchName()

	var branchGitHash string
	for _, branch := range allBranches {
		if branch[strings.LastIndex(branch, "/")+1:] == branchName {
			branchGitHash = branch
		}
	}

	if branchExists {
		branchGitHash = strings.Fields(branchGitHash)[0]
		return branchGitHash[0:7]
	} else {
		return ""
	}

}

func GetMasterGitHash() string {
	allRepoBranches := GetAllRepoBranches()
	masterHashStringArray := strings.Fields(allRepoBranches[0])
	masterHash := masterHashStringArray[0][0:7]
	return masterHash
}

func GetBranchFlagAndName() (bool, string) {
	args := os.Args[1:]

	//if branchName is empty string, strings.Contain does not work
	//this is a workaround
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

package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// CloneRepo grabs url and clones the repo/copies current dir
func CloneRepo() {

	//clone repo with url from args
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	//if builder cmd, copy current dir to hidden instead of clone
	if os.Getenv("BUILDER_COMMAND") == "true" {
		//pwd
		path, err := os.Getwd()
		if err != nil {
			BuilderLog.Errorf("failed to get repository", err)
		}
		fmt.Println(path)
		exec.Command("cp", "-a", path+"/.", hiddenDir).Run()
	} else {
		repo := GetRepoURL()
		//if config cmd clone to temp dir on first go
		if hiddenDir == "" {
			errDir := os.Mkdir("./tempRepo", 0755)
			if errDir != nil {
				fmt.Println(errDir)
			}

			cmd := exec.Command("git", "clone", repo, "./tempRepo")
			cmd.Run()
		} else {
			bFlagExists, branchExists, branchName := bFlagAndBranchExists()

			if bFlagExists {
				if branchExists {
					cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, hiddenDir)
					cmd.Run()
					BuilderLog.Infof("git clone", cmd)
				} else {
					BuilderLog.Fatal("Branch does not exists")
				}
			} else {
				cmd := exec.Command("git", "clone", repo, hiddenDir)
				cmd.Run()
				BuilderLog.Infof("git clone", cmd)
			}
		}
	}
}

func bFlagAndBranchExists() (bool, bool, string) {
	//if init cmd, clone to hidden dir
	bFlagExists, branchName := CloneBranch()

	//check if branch exist
	branches, _, _, _ := GitHashAndName()
	branchExists, _ := BranchNameExists(branches)

	return bFlagExists, branchExists, branchName
}

func CloneBranch() (bool, string) {
	args := os.Args[1:]

	//if branch is empty string strings.Contain does not work, function found in metadata
	branchName := "%$F"
	branchFlag := false

	repoBranch, present := os.LookupEnv("REPO_BRANCH")
	if present && (repoBranch != "") && !(contains(args, "-b") || contains(args, "--branch")) {
		//convert val interface{} to string to be set as env var
		branchName = repoBranch
		branchFlag = true
	} else {
		for i, v := range args {
			if v == "-b" || v == "--branch" {
				if len(args) <= i+1 {
					BuilderLog.Fatal("No Branch Name Provided")

				} else {
					branchName = args[i+1]
					branchFlag = true
				}
			}
		}
	}

	return branchFlag, branchName
}

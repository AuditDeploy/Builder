package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

			// Get branch name
			os.Setenv("REPO_BRANCH_NAME", GetBranchName())
		} else {
			bFlagExists, branchExists, branchName := bFlagAndBranchExists()

			if bFlagExists {
				if branchExists {
					cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, hiddenDir)
					cmd.Run()
					BuilderLog.Infof("git clone", cmd)

					// Get branch name
					os.Setenv("REPO_BRANCH_NAME", branchName)
				} else {
					BuilderLog.Fatal("Branch does not exists")
				}
			} else if branchExists { // Repo branch given in builder.yaml not by -b flag
				cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, hiddenDir)
				cmd.Run()
				BuilderLog.Infof("git clone", cmd)

				// Get branch name
				os.Setenv("REPO_BRANCH_NAME", branchName)
			} else {
				cmd := exec.Command("git", "clone", repo, hiddenDir)
				cmd.Run()
				BuilderLog.Infof("git clone", cmd)

				// Get branch name
				os.Setenv("REPO_BRANCH_NAME", GetBranchName())
			}
		}
	}
}

func bFlagAndBranchExists() (bool, bool, string) {
	var bFlagExists, branchExists = false, false
	args := os.Args[1:]

	branchName := os.Getenv("REPO_BRANCH")

	for i, v := range args {
		if v == "-b" || v == "--branch" {
			if len(args) <= i+1 {
				BuilderLog.Fatal("No Branch Name Provided")

			} else {
				branchName = args[i+1]
				bFlagExists = true
			}
		}
	}

	if branchName != "" {
		branchExists = true
	}

	return bFlagExists, branchExists, branchName
}

func GetBranchName() string {
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmd.Dir = hiddenDir
	branch, _ := branchCmd.Output()

	return strings.TrimSuffix(string(branch), "\n")
}

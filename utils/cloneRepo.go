package utils

import (
	"Builder/spinner"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneRepo grabs url and clones the repo
func CloneRepo(to string) {
	// Get absolute path if a relative path was given
	toPath, _ := filepath.Abs(to)

	repo := GetRepoURL()

	// Stat path, if it doesn't exist, create it
	_, err := os.Stat(toPath)
	if err != nil {
		errDir := os.MkdirAll(toPath, 0755)
		if errDir != nil {
			spinner.LogMessage("Could not create new repo directory: "+errDir.Error(), "fatal")
		}

		cmd := exec.Command("git", "clone", repo, toPath)
		cmd.Run()

		// Get branch name
		os.Setenv("REPO_BRANCH_NAME", GetBranchName(toPath))
	} else {
		bFlagExists, branchExists, branchName := bFlagAndBranchExists()

		if bFlagExists {
			if branchExists {
				cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, toPath)
				cmd.Run()
				spinner.LogMessage("git clone -b "+branchName+" --single-branch "+repo, "info")

				// Get branch name
				os.Setenv("REPO_BRANCH_NAME", branchName)
			} else {
				spinner.LogMessage("Branch does not exists", "fatal")
			}
		} else if branchExists { // Repo branch given in builder.yaml not by -b flag
			cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repo, toPath)
			cmd.Run()
			spinner.LogMessage("git clone -b "+branchName+" --single-branch "+repo, "info")

			// Get branch name
			os.Setenv("REPO_BRANCH_NAME", branchName)
		} else {
			cmd := exec.Command("git", "clone", repo, toPath)
			cmd.Run()
			spinner.LogMessage("git clone "+repo, "info")

			// Get branch name
			os.Setenv("REPO_BRANCH_NAME", GetBranchName(toPath))
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
				spinner.LogMessage("No Branch Name Provided", "fatal")

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

func GetBranchName(path string) string {
	branchCmd := exec.Command("git", "branch")
	branchCmd.Dir = path
	branch, _ := branchCmd.Output()

	formatName := strings.TrimSuffix(string(branch), "\n")

	return formatName[2:]
}

package utils

import (
	"Builder/spinner"
	"os"
	"os/exec"
)

func GetRepoURL() string {
	args := os.Args[1:]

	//grab URL first
	var repo string
	for i, v := range args {
		if v == "init" || v == "config" {
			if len(args) <= i+1 {
				// logger.ErrorLogger.Println("No Repo Url Provided")
				spinner.LogMessage("No Repo Url Provided", "fatal")

			} else {
				repo = args[i+1]
			}
		}
	}
	if repo == "" {
		// Get repo name from git config file
		out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		if err != nil {
			spinner.LogMessage("Can't get repo url from .git/config file: "+err.Error(), "fatal")
		}

		repo = string(out[:len(out)-1]) // Get rid of last char because it is a newline char

		if repo == "" {
			spinner.LogMessage("Can't get repo url from .git/config file", "fatal")
		}

	}
	return repo
}

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
			spinner.LogMessage("Can't find git URL.  Please provide it in the builder.yaml", "info")
			return ""
		}

		repo = string(out[:len(out)-1]) // Get rid of last char because it is a newline char

		if repo == "" {
			if os.Getenv("GIT_URL") != "" {
				repo = os.Getenv("GIT_URL")
			} else {
				spinner.LogMessage("Can't find git URL.  Please provide it in the builder.yaml", "info")
			}
		}

	}
	return repo
}

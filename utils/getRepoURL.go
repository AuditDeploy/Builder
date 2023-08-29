package utils

import (
	"Builder/spinner"
	"os"
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
		// logger.ErrorLogger.Println("No Repo Url Provided")
		spinner.LogMessage("No Repo Url Provided", "fatal")
	}
	return repo
}

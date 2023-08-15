package utils

import (
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
				BuilderLog.Fatal("No Repo Url Provided")

			} else {
				repo = args[i+1]
			}
		}
	}
	if repo == "" {
		// logger.ErrorLogger.Println("No Repo Url Provided")
		BuilderLog.Fatal("No Repo Url Provided")
	}
	return repo
}

package utils

import (
	"fmt"
	"os"
)

func GetRepoURL() string {
	args := os.Args[1:]

	//grab URL first
	var repo string
	for i, v := range args {
		if v == "--repo" || v == "-r" {
			if len(args) <= i+1 {
				fmt.Println("No Repo Url Provided")
				os.Exit(1)

			} else {
				repo = args[i+1]
			}
		}
	}
	if repo == "" {
		fmt.Println("No Repo Url Provided")
		os.Exit(1)
	}
	return repo
}

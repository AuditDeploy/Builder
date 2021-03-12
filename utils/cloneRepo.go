package utils

import (
	"fmt"
	"os"
	"os/exec"
)

//CloneRepo grabs url
func CloneRepo() {

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

	//clone repo with url from args
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	// enter parent name/hidden dir
	cmd := exec.Command("git", "clone", repo, hiddenDir)
	cmd.Run()
}


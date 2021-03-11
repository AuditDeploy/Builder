package utils

import (
	"fmt"
	"os"
	"os/exec"
)

//CheckArgs is...
func CheckArgs() {
	var repo string
	cArgs := os.Args[1:]

	//check repo flag syntax
	for i, v := range cArgs {
		if v == "--repo" || v == "-r" {
			if len(cArgs) <= i+1 {
				fmt.Println("No Repo Url Provided")
				os.Exit(1)

			} else {
				 repo = cArgs[i+1]
			}
		}
	}

	//if flag present, but no url
	if repo == "" {
		fmt.Println("No Repo Url Provided")
		os.Exit(1)
	}

	//check to see if repo exists
	//returns the exit status in err
	_, err := 	exec.Command("git", "ls-remote", repo, "-q").Output()
	if (err != nil) { 
		fmt.Println("Repo Provided Does Not Exists")
		os.Exit(1)
	} 
}
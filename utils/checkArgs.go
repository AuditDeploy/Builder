package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//CheckArgs is...
func CheckArgs() {
	//Repo
	repo := GetRepoURL()

	cArgs := os.Args[1:]

	//check for help flag
	for _, v := range cArgs {
		if v == "--help" || v == "-h" {
			fmt.Println("**Info about Builder**")
			os.Exit(1)
		}
	}

	//if flag present, but no url
	if repo == "" {
		log.Fatal("No Repo Url Provided")
	}

	//check to see if repo exists
	//git ls-remote lists refs/heads & tags of a repo, if none exists, exit status thrown
	//returns the exit status in err
	_, err := exec.Command("git", "ls-remote", repo, "-q").Output()
	if err != nil {
		log.Fatal("Repo Provided Does Not Exists")
	}

	//check if artifact path is passed in 
	var artifactPath string
	for i, v := range cArgs { 
		if v == "--path" || v == "-p" {
			if len(cArgs) <= i+1 {
				log.Fatal("No Path Provided")

			} else {
				artifactPath = cArgs[i+1]
			}
		}
	}
	val, present := os.LookupEnv("BUILDER_OUTPUT_PATH")
	if !present {
		os.Setenv("BUILDER_OUTPUT_PATH", artifactPath)
	} else {
		fmt.Println("BUILDER_OUTPUT_PATH", val)
	}
}

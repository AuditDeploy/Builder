package utils

import (
	"Builder/spinner"
	"fmt"
	"os"
	"os/exec"
)

// CheckArgs is...
func CheckArgs() {
	//Repo
	repo := GetRepoURL()
	cArgs := os.Args[1:]
	//if flag present, but no url
	if repo == "" {
		spinner.LogMessage("No Repo Url Provided", "fatal")
	}

	//check to see if repo exists
	//git ls-remote lists refs/heads & tags of a repo, if none exists, exit status thrown
	//returns the exit status in err
	_, err := exec.Command("git", "ls-remote", repo, "-q").Output()
	if err != nil {
		spinner.LogMessage("Provided repository does not exist", "fatal")
	}

	//check if artifact path is passed in
	var artifactPath string
	for i, v := range cArgs {
		if v == "--output" || v == "-o" {
			if len(cArgs) <= i+1 {
				spinner.LogMessage("No Output Path Provided", "fatal")

			} else {
				artifactPath = cArgs[i+1]
				val, present := os.LookupEnv("BUILDER_OUTPUT_PATH")
				if !present {
					os.Setenv("BUILDER_OUTPUT_PATH", artifactPath)
				} else {
					fmt.Println("BUILDER_OUTPUT_PATH", val)
					fmt.Println("Output Path already present")
					spinner.LogMessage("Output path already present", "error")
				}
			}
		}
		if v == "--compress" || v == "-z" || v == "-C" {
			os.Setenv("ARTIFACT_ZIP_ENABLED", "true")
		}
		if v == "--hidden" || v == "-H" {
			os.Setenv("HIDDEN_DIR_ENABLED", "true")
		}
	}
}

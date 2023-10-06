package utils

import (
	"Builder/spinner"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
					// If init or config command used and relative path provided,
					// make sure path is relative to project folder
					for _, a := range cArgs {
						if a == "config" || a == "init" {
							if !filepath.IsAbs(artifactPath) {
								if strings.HasPrefix(artifactPath, "./") {
									artifactPath = strings.Replace(artifactPath, "./", "./"+GetName(), 1)
								} else if strings.HasPrefix(artifactPath, "../") {
									artifactPath = strings.Replace(artifactPath, "../", "./", 1)
								}
							}
						}
					}

					// Make sure output path is absolute
					outputPath, err := filepath.Abs(artifactPath)
					if err != nil {
						spinner.LogMessage("Could not resolve outputpath", "fatal")
					}

					os.Setenv("BUILDER_OUTPUT_PATH", outputPath)
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

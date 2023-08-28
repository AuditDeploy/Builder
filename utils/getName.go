package utils

import (
	"Builder/spinner"
	"os"
	"strings"
)

// GetName does ...
func GetName() string {
	var name string
	args := os.Args[1:]

	dirName, present := os.LookupEnv("BUILDER_DIR_NAME")
	if present && (dirName != "") && !(contains(args, "-n") || contains(args, "--name")) {
		//convert val interface{} to string to be set as env var
		name = dirName
	} else {
		//if args contains name flag, use name
		if contains(args, "-n") || contains(args, "--name") {
			for i, v := range args {
				if v == "--name" || v == "-n" {
					if len(args) <= i+1 {
						spinner.LogMessage("Please provide a name", "fatal")
					} else {
						if specialChar(args[i+1]) {
							spinner.LogMessage("Special Characters Not Allowed In Names", "fatal")
						}
						name = args[i+1]
					}
				}
			}
		} else if os.Getenv("BUILDER_COMMAND") == "true" {
			//use current dir name if no --name flag and using builder cmd
			path, err := os.Getwd()
			if err != nil {
				spinner.LogMessage("error getting builder command directory", "error")
			}
			name = path[strings.LastIndex(path, "/")+1:]

		} else {
			//if init or config and no --name flag, use repo name
			repoURL := os.Args[2]
			name = repoURL[strings.LastIndex(repoURL, "/")+1:]

			//if .git still in the name, remove it
			if strings.HasSuffix(name, ".git") {
				name = strings.TrimSuffix(name, ".git")
			}
		}
		os.Setenv("BUILDER_DIR_NAME", name)
	}
	return name
}

func specialChar(str string) bool {
	hasSpecialCharacter := false

	f := func(r rune) bool {
		return r < 'A' || r > 'z'
	}

	if strings.IndexFunc(str, f) != -1 {
		hasSpecialCharacter = true
	}

	return hasSpecialCharacter
}

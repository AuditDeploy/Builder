package utils

import (
	"log"
	"os"
	"strings"
)

//GetName does ...
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
						log.Fatal("Please provide a name")
					} else {
						if specialChar(args[i+1]) {
							log.Fatal("Special Characters Not Allowed In Names")
						}
						name = args[i+1]
					}
				}
			}
		} else if os.Getenv("BUILDER_COMMAND") == "true" {
			//use current dir name if no --name flag and using builder cmd
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			name = path[strings.LastIndex(path, "/")+1:]

		} else {
			//if init or config and no --name flag, use repo name
			repoURL := os.Args[2]
			name = repoURL[strings.LastIndex(repoURL, "/")+1:]
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

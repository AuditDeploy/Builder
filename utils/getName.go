package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

//GetName does ...
func GetName() string {
	var name string
	args := os.Args[1:]

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
		} else if os.Getenv("BUILDER_COMMAND") == "true" {
			fmt.Print("No name detected. Use '-n' to define one.")
		} else {
			repoURL := os.Args[2]
			name = repoURL[strings.LastIndex(repoURL, "/")+1:]
		}
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

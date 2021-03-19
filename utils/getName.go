package utils

import (
	"log"
	"os"
	"strings"
)

//GetName does ...
func GetName(cArgs []string) string {
	repoURL := os.Args[3]

	name := repoURL[strings.LastIndex(repoURL, "/")+1:]

	for i, v := range cArgs {
		if v == "--name" || v == "-n" {
			if len(cArgs) <= i+1 {
				log.Fatal("Please provide a name")
			} else {
				if specialChar(cArgs[i+1]) {
					log.Fatal("Special Characters Not Allowed In Names")
				}
				name = cArgs[i+1]
			}
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

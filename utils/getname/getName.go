package getname

import (
	"log"
	"os"
	"strings"
)

//GetName does ...
func GetName(cArgs []string) string {
	repoURL := os.Args[2]

	name := repoURL[strings.LastIndex(repoURL, "/")+1:]

	for i, v := range cArgs {
		if v == "--name" || v == "-n" {
			if len(cArgs) <= i+1 {
				log.Fatal("Please provide a name")
			} else {
				name = cArgs[i+1]
			}
		}
	}

	return name
}

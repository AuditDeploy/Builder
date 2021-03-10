package repo

import (
	"fmt"
	"os"
	"os/exec"
)

//GetURL grabs url
func GetURL(cArgs []string) {
	var repo string
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
	if repo == "" {
		fmt.Println("No Repo Url Provided")
		os.Exit(1)
	}
	cloneRepo(repo)
}

func cloneRepo(repo string) {
	fmt.Println(repo)

	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	// enter parent name/hidden dir
	cmd := exec.Command("git", "clone", repo, hiddenDir)
	cmd.Run()

	//make contents read-only
	cmd2 := exec.Command("chmod", "-R", "0444", hiddenDir)
	cmd2.Run()  
}

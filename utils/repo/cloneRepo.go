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

	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	// enter parent name/workspace dir
	cmd := exec.Command("git", "clone", repo, workspaceDir)
	cmd.Run()
}

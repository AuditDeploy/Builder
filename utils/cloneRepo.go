package utils

import (
	"os"
	"os/exec"
)

//CloneRepo grabs url
func CloneRepo() {

	repo := GetRepoURL()

	//clone repo with url from args
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	// enter parent name/hidden dir
	cmd := exec.Command("git", "clone", repo, hiddenDir)
	cmd.Run()
}

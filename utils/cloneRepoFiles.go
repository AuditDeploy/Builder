package utils

import (
	"Builder/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneRepo grabs url and clones the repo/copies current dir
func CloneRepoFiles(from string, to string) {
	//Get absolute paths in case relative paths are given
	fromPath, _ := filepath.Abs(from)
	toPath, _ := filepath.Abs(to)

	files, err := os.ReadDir(fromPath)
	if err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	// Copy all but Builder created dir
	if os.Getenv("BUILDER_BUILDS_DIR") != "" { // A different folder to store builds is provided
		for _, file := range files {
			if file.Name() != os.Getenv("BUILDER_BUILDS_DIR") {
				//copy files to given path
				cmd := exec.Command("cp", "-r", file.Name(), toPath)
				cmd.Dir = fromPath
				cmd.Run()
			}
		}
	} else { // Builds are stored in default named 'builder' folder
		for _, file := range files {
			if file.Name() != "builder" {
				//copy files to given path
				cmd := exec.Command("cp", "-r", file.Name(), toPath)
				cmd.Dir = fromPath
				cmd.Run()
			}
		}
	}
}

package utils

import (
	"Builder/spinner"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

// CloneRepo grabs url and clones the repo/copies current dir
func CloneRepoFiles(from string, to string) {
	//Get absolute paths in case relative paths are given
	fromPath, _ := filepath.Abs(from)
	toPath, _ := filepath.Abs(to)

	// Copy all but Builder created dir
	if os.Getenv("BUILDER_BUILDS_DIR") != "" { // A different folder to store builds is provided
		opt := cp.Options{
			Skip: func(info os.FileInfo, src, dest string) (bool, error) {
				return info.Name() == os.Getenv("BUILDER_BUILDS_DIR"), nil
			},
		}
		err := cp.Copy(fromPath, toPath, opt)
		if err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}
	} else { // Builds are stored in default named 'builder' folder
		// copy files to given path
		opt := cp.Options{
			Skip: func(info os.FileInfo, src, dest string) (bool, error) {
				return info.Name() == "builder", nil
			},
		}
		err := cp.Copy(fromPath, toPath, opt)
		if err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}
	}
}

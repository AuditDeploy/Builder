package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/utils/log"
)

func Init() {
	// Check argument syntax, exit if incorrect
	utils.CheckArgs()

	// Make dirs
	directory.MakeDirs()
	log.Info("Directories successfully created.")

	// Clone repo into hidden
	utils.CloneRepo()
	log.Info("Repo cloned successfully.")

	// Compile logic to derive project type
	derive.ProjectType()

	// Get build metadata (deprecated, func moved inside compiler)
	// utils.Metadata()
	log.Info("Metadata created successfully.")

	// Check for Dockerfile, then build image
	utils.Docker()

	// Makes hidden dir read-only
	utils.MakeHidden()
	log.Info("Hidden Dir is now read-only.")
}

package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// make dirs
	directory.MakeDirs()
	BuilderLog.Info("Directories successfully created.")

	// clone repo into hidden
	utils.CloneRepo()
	BuilderLog.Info("Repo cloned successfully.")

	// compile logic to derive project type
	derive.ProjectType()

	//Get build metadata (deprecated, func moved inside compiler)
	// utils.Metadata()
	BuilderLog.Info("Metadata created successfully.")

	// Store build metadata to hidden builder dir
	utils.StoreBuildMetadataLocally()

	//Check for Dockerfile, then build image
	utils.Docker()

	//makes hidden dir read-only
	utils.MakeHidden()
	BuilderLog.Info("Hidden Dir is now read-only.")
}

package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// Start loading spinner
	spinner.Spinner.Start()

	// make dirs
	directory.MakeDirs()
	spinner.LogMessage("Directories successfully created.", "info")

	// clone repo into hidden
	utils.CloneRepo()
	spinner.LogMessage("Repo cloned successfully.", "info")

	// compile logic to derive project type
	derive.ProjectType()

	//Get build metadata (deprecated, func moved inside compiler)
	// utils.Metadata()
	spinner.LogMessage("Metadata created successfully.", "info")

	// Store build metadata to hidden builder dir
	utils.StoreBuildMetadataLocally()

	//Check for Dockerfile, then build image
	utils.Docker()

	//makes hidden dir read-only
	utils.MakeHidden()
	spinner.LogMessage("Hidden Dir is now read-only.", "info")

	// Stop loading spinner
	spinner.Spinner.Stop()
}

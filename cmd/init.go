package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"

	"os"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// Start loading spinner
	spinner.Spinner.Start()

	// clone repo into folder named after project name and keep track of its path
	projectName := utils.GetName()
	os.Setenv("BUILDER_REPO_DIR", "./"+projectName)
	utils.CloneRepo("./" + projectName)
	spinner.LogMessage("Repo cloned successfully.", "info")

	// make dirs
	directory.MakeDirs()
	spinner.LogMessage("Directories successfully created.", "info")

	// copy repo files we just cloned to hidden dir
	repoDir := os.Getenv("BUILDER_REPO_DIR")
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	utils.CloneRepoFiles(repoDir, hiddenDir)

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

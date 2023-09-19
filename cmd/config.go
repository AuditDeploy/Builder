package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/yaml"

	"os"
)

func Config() {
	//check args normally,
	utils.CheckArgs()

	// Start loading spinner
	spinner.Spinner.Start()

	//clone repo into temp dir to pull builder.yaml info
	utils.CloneRepo("./tempRepo")

	//set yaml info as env vars
	yaml.YamlParser("./tempRepo/builder.yaml")

	// clone repo into folder named after project name and keep track of its path
	projectName := utils.GetName()
	if os.Getenv("BUILDER_DIR_PATH") != "" { // User wants this repo built elsewhere
		utils.CloneRepo(os.Getenv("BUILDER_DIR_PATH") + "/" + projectName)
		os.Setenv("BUILDER_REPO_DIR", os.Getenv("BUILDER_DIR_PATH")+"/"+projectName)
	} else {
		utils.CloneRepo("./" + projectName)
		os.Setenv("BUILDER_REPO_DIR", "./"+projectName)
	}
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

package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/yaml"
)

func Config() {
	//check args normally,
	utils.CheckArgs()

	//clone repo into temp dir to pull builder.yaml info
	utils.CloneRepo()

	//set yaml info as env vars
	yaml.YamlParser("./tempRepo/builder.yaml")

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

package cmd

import (
	"Builder/directory"
	"Builder/logger"
	"Builder/utils"
	"Builder/yaml"
)

func Config() {
	//check args normally,
	utils.CheckArgs()
	
  //clone repo into temp dir to pull builder.yaml info 
	utils.CloneRepo()

	//set yaml info as env vars
	yaml.YamlParser()

	// make dirs
	directory.MakeDirs()
	logger.InfoLogger.Println("Directories successfully created.")

	// clone repo into hidden
	utils.CloneRepo()
	logger.InfoLogger.Println("Repo cloned successfully.")

	// compile logic to derive project type
	utils.ProjectType()

	//Get build metadata
	utils.Metadata()
	logger.InfoLogger.Println("Metadata created successfully.")

	//makes hidden dir read-only
	utils.MakeHidden()
	logger.InfoLogger.Println("Hidden Dir is now read-only.")

	//creates global logs dir
	logger.GlobalLogs()
	// delete temp dir
}
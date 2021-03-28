package cmd

import (
	directory "Builder/directory"
	"Builder/logger"
	utils "Builder/utils"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// make dirs
	directory.MakeParentDir()
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
}

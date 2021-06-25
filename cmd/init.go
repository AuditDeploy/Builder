package cmd

import (
	"builder/derive"
	"builder/directory"
	"builder/logger"
	"builder/utils"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// make dirs
	directory.MakeDirs()
	logger.InfoLogger.Println("Directories successfully created.")

	// clone repo into hidden
	utils.CloneRepo()
	logger.InfoLogger.Println("Repo cloned successfully.")

	// compile logic to derive project type
	derive.ProjectType()

	//Get build metadata
	// utils.Metadata()
	logger.InfoLogger.Println("Metadata created successfully.")

	//makes hidden dir read-only
	utils.MakeHidden()
	logger.InfoLogger.Println("Hidden Dir is now read-only.")

	//creates global logs dir
	logger.GlobalLogs()
}

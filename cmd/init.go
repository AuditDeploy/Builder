package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/utils/log"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// make dirs
	directory.MakeDirs()
	log.Info("Directories successfully created.")

	// clone repo into hidden
	utils.CloneRepo()
	log.Info("Repo cloned successfully.")

	// compile logic to derive project type
	derive.ProjectType()

	//Get build metadata (deprecated, func moved inside compiler)
	// utils.Metadata()
	log.Info("Metadata created successfully.")

	// Change parent dir of build to include build start time
	directory.UpdateParentDirName()

	//Check for Dockerfile, then build image
	utils.Docker()

	//makes hidden dir read-only
	utils.MakeHidden()
	log.Info("Hidden Dir is now read-only.")
}

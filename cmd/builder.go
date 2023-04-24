package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/logger"
	"Builder/utils"
	"Builder/yaml"
	"log"
	"os"
	"os/exec"
)

func Builder() {
	os.Setenv("BUILDER_COMMAND", "true")
	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		exec.Command("git", "pull").Run()

		//pareses builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		//append logs
		//logger.CreateLogs(os.Getenv("BUILDER_LOGS_DIR"))
		directory.MakeDirs()
		logger.InfoLogger.Println("Directories successfully created.")

		// clone repo into hidden
		utils.CloneRepo()
		logger.InfoLogger.Println("Repo cloned successfully.")

		//creates a new artifact
		derive.ProjectType()

		//Get build metadata (deprecated, func moved inside compiler)
		logger.InfoLogger.Println("Metadata created successfully.")

		//Check for Dockerfile, then build image
		utils.Docker()

		//makes hidden dir read-only
		utils.MakeHidden()
		logger.InfoLogger.Println("Hidden Dir is now read-only.")

		//creates global logs dir
		logger.GlobalLogs()
		// delete temp dir
	} else {
		log.Fatal("bulder.yaml file not found. cd into it's location.")
	}
}

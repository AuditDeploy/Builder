package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
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
		log.Info("Directories successfully created.")

		// clone repo into hidden
		utils.CloneRepo()
		log.Info("Repo cloned successfully.")

		//creates a new artifact
		derive.ProjectType()

		//Get build metadata (deprecated, func moved inside compiler)
		log.Info("Metadata created successfully.")

		//Check for Dockerfile, then build image
		utils.Docker()

		//makes hidden dir read-only
		utils.MakeHidden()
		log.Info("Hidden Dir is now read-only.")

	} else {
		utils.Help()
	}
}

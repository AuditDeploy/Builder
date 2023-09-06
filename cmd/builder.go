package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/yaml"
	"os"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()

func Builder() {
	os.Setenv("BUILDER_COMMAND", "true")
	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {

		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		// Create directories
		directory.MakeDirs()
		BuilderLog.Info("Directories successfully created.")

		// clone repo into hidden
		//utils.CloneRepo()
		//BuilderLog.Info("Repo cloned successfully.")

		//creates a new artifact
		derive.ProjectType()

		//Get build metadata (deprecated, func moved inside compiler)
		BuilderLog.Info("Metadata created successfully.")

		//Check for Dockerfile, then build image
		utils.Docker()

		//makes hidden dir read-only
		utils.MakeHidden()
		BuilderLog.Info("Hidden Dir is now read-only.")

	} else {
		utils.Help()
	}
}

package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/utils"
	"Builder/yaml"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()

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
		BuilderLog.Info("Directories successfully created.")

		// clone repo into hidden
		utils.CloneRepo()
		BuilderLog.Info("Repo cloned successfully.")

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

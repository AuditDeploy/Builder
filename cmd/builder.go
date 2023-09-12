package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/yaml"
	"os"
)

func Builder() {
	os.Setenv("BUILDER_COMMAND", "true")
	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		// Start loading spinner
		spinner.Spinner.Start()

		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		// Create directories
		directory.MakeDirs()
		spinner.LogMessage("Directories successfully created.", "info")

		// clone files from current dir into hidden
		currentDir, _ := os.Getwd()
		hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
		utils.CloneRepoFiles(currentDir, hiddenDir)
		spinner.LogMessage("Files copied to hidden dir successfully.", "info")

		//creates a new artifact
		derive.ProjectType()

		//Get build metadata (deprecated, func moved inside compiler)
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
	} else {
		utils.Help()
	}
}

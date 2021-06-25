package cmd

import (
	"builder/derive"
	"builder/logger"
	"builder/yaml"
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

		//parses builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		//append logs
		logger.CreateLogs(os.Getenv("BUILDER_LOGS_DIR"))

		//creates a new artifact
		derive.ProjectType()

	} else {
		log.Fatal("bulder.yaml file not found. cd into it's location.")
	}
}

package main

import (
	"Builder/cmd"
	"Builder/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) > 1 {
		builderCommand := os.Args[1]

		if builderCommand == "init" {
			cmd.Init()
		} else if builderCommand == "config" {
			cmd.Config()
		} else {
			fmt.Println("expected command: 'init' or 'config'")
		}
	} else {
		builder()
	}
}

func builder() {
	os.Setenv("BUILDER_CMD", "true")

	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		exec.Command("git", "pull").Run()

		//pareses builder.yaml
		utils.YamlParser(path + "/" + "builder.yaml")

		// projectPath := os.Getenv("BUILDER_DIR_PATH")
		// projectType := os.Getenv("BUILDER_PROJECT_TYPE")
		// buildTool := os.Getenv("BUILDER_BUILD_TOOL")
		// buildFile := os.Getenv("BUILDER_BUILD_FILE")
		// buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

		//create binary file
		// utils.ProjectType()

	} else {
		log.Fatal("bulder.yaml file not found or cd into workspace")
	}

}

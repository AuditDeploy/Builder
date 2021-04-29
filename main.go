package main

import (
	"Builder/cmd"
	"Builder/derive"
	"Builder/logger"
	"Builder/yaml"
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
	os.Setenv("BUILDER_COMMAND", "true")
	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		exec.Command("git", "pull").Run()

		//pareses builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		fmt.Println("before logs")
		//append logs
		logger.CreateLogs(os.Getenv("BUILDER_LOGS_DIR"))

		fmt.Println("before derive")
		//run derive 
		//creates a new binary
		derive.ProjectType()

	} else {
		log.Fatal("bulder.yaml file not found or cd into workspace")
	}
}

package main

import (
	"Builder/cmd"
	"Builder/yaml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	path, _ := os.Getwd()

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		exec.Command("git", "pull").Run()

		//pareses builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		//creates a new binary
		buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
		buildCmdArray := strings.Fields(buildCmd)
		exec.Command(buildCmdArray[0], buildCmdArray[1:]...).Run()

	} else {
		log.Fatal("bulder.yaml file not found or cd into workspace")
	}
}

package utils

import (
	"Builder/logger"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Docker creates image from dockerfile and pushes to dockerhub
func Docker() {
	dockerFlag := CheckDockerFlag()

	//if -d flag exists, build image
	if dockerFlag {
		//DETERMINE CMD
		var cmd *exec.Cmd
		dockerCmd := os.Getenv("BUILDER_DOCKER_CMD")
		//if dockerCmd doesn't exist use default
		if dockerCmd == "" {
			name := GetName()
			imageName := fmt.Sprintf("builder/%s", name)
			unixTime := os.Getenv("BUILDER_TIMESTAMP")
			cmd = exec.Command("docker", "build", ".", "-t", imageName+"-"+unixTime)
		} else {
			//else use defined dockerCmd
			dockerCmdArray := strings.Fields(dockerCmd)
			cmd = exec.Command(dockerCmdArray[0], dockerCmdArray[1:]...)
		}

		//DETERMINE PATH
		//determine projectType to top level Dockerfile path
		compType := []string{"go", "c#", "java"}
		nonCompType := []string{"node", "npm", "python", "ruby"}
		workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
		projectType := os.Getenv("BUILDER_PROJECT_TYPE")
		if contains(compType, projectType) {
			cmd.Dir = workspaceDir
		} else if contains(nonCompType, projectType) {
			cmd.Dir = workspaceDir + "/temp/"
		} else {
			log.Fatal("Please define your projectType in builder.yaml")
		}

		//RUN DOCKER BUILD
		logger.InfoLogger.Println(cmd)
		err := cmd.Run()
		if err != nil {
			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			logger.ErrorLogger.Println("Docker build failed.")
			fmt.Println("out:", outb.String(), "err:", errb.String())
			log.Fatal(err)
		}

		//RUN DOCKER PUSH
	}
}

//CheckDockerFlag for docker flag
func CheckDockerFlag() bool {
	var exists bool
	cArgs := os.Args[1:]
	for _, v := range cArgs {
		if v == "--docker" || v == "-d" {
			fmt.Println("Building Docker Image... 🐳")
			exists = true
		} else {
			exists = false
		}
	}
	return exists
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

/****
TODO:
- check project type to determine cmd path***
- allow user pass creds on cmd line (probably don't want it to live in builder.yaml)
	- if not present, run anyway with whatever log in on users machine
- allow user to input command in yaml (split string by spaces)***
	- if not present, run default with name of project

THEN:
- Push to dockerhub with credentials
****/

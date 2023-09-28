package utils

import (
	"Builder/spinner"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Docker creates image from dockerfile and pushes to dockerhub
func Docker() {
	dockerFlag := CheckDockerFlag()

	//if -D flag exists, build image
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
		compType := []string{"go", "rust", "c#", "java"}
		nonCompType := []string{"node", "npm", "python", "ruby"}
		workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
		projectType := os.Getenv("BUILDER_PROJECT_TYPE")
		if contains(compType, projectType) {
			cmd.Dir = workspaceDir
		} else if contains(nonCompType, projectType) {
			cmd.Dir = workspaceDir + "/temp/"
		} else {
			spinner.LogMessage("Please define your projectType in builder.yaml", "fatal")
		}

		//RUN DOCKER BUILD
		spinner.LogMessage("running command: "+cmd.String(), "info")
		err := cmd.Run()
		if err != nil {
			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			fmt.Println("out:", outb.String(), "err:", errb.String())
			spinner.LogMessage("docker build failed: "+err.Error(), "fatal")
		}

		//RUN DOCKER PUSH
	}
}

// CheckDockerFlag for docker flag
func CheckDockerFlag() bool {
	var exists bool
	cArgs := os.Args[1:]
	for _, v := range cArgs {
		if v == "--docker" || v == "-D" {
			spinner.LogMessage("Building docker image 🐳", "info")
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

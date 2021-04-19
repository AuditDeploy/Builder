package compile

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CSharp(filePath string) {
	fmt.Println(filePath)

	//define dir path for command to run in
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included already, else add curren dir + local path 
	if (configPath != "") {
		// ex: C:/Users/Name/Projects/helloworld_19293/workspace/dir
		fullPath = filePath
	} else {
		path, _ := os.Getwd()
		//combine local path to newly created tempWorkspace, 
		//gets rid of "." in path name
		// ex: C:/Users/Name/Projects + /helloworld_19293/workspace/dir
		fullPath = path + filePath[strings.Index(filePath, ".")+1:]
	}

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	var cmd *exec.Cmd
	if buildTool == "dotnet" {
		fmt.Println(buildTool)
		cmd = exec.Command("dotnet", "build")
		cmd.Dir = fullPath       // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("dotnet", "build")
		cmd.Dir = fullPath       // or whatever directory it's in
	}

	fmt.Println(cmd.Dir)

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("C# project failed to compile.")
		log.Fatal(err)
	}

	logger.InfoLogger.Println("C# project compiled successfully.")
}

package compile

import (
	"Builder/logger"
	"Builder/yaml"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CSharp(filePath string) {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "c#")
	}

	//define dir path for command to run in
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included already, else add curren dir + local path
	if configPath != "" {
		// ex: C:/Users/Name/Projects/helloworld_19293/workspace/dir
		fullPath = filePath
	} else {
		path, _ := os.Getwd()
		//combine local path to newly created tempWorkspace,
		//gets rid of "." in path name
		// ex: C:/Users/Name/Projects + /helloworld_19293/workspace/dir
		fullPath = path + filePath[strings.Index(filePath, ".")+1:]
		os.Setenv("BUILDER_DIR_PATH", path)
	}

	//install dependencies/build,
	// if yaml build type exists install accordingly, if buildCmd exists,
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		cmd = exec.Command(buildCmd)
	} else if buildTool == "dotnet" {
		cmd = exec.Command("dotnet", "build", fullPath)
		cmd.Dir = fullPath // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("dotnet", "build", fullPath)
		// cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_TOOL", "dotnet")
		os.Setenv("BUILDER_BUILD_COMMAND", "dotnet build "+fullPath)
		os.Setenv("BUILDER_BUILD_FILE", fullPath[strings.LastIndex(fullPath, "/")+1:])
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("C# project failed to compile.")
		log.Fatal(err)
	}

	fullPath = fullPath[:strings.LastIndex(fullPath, "/")+1]
	yaml.CreateBuilderYaml(fullPath)

	// artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	// if (artifactPath != "") {
	// 	exec.Command("cp", "-a", fullPath+"/main.exe", artifactPath).Run()
	// }

	logger.InfoLogger.Println("C# project compiled successfully.")
}

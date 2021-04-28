// takes in code as arg from go
//run go build on code given

package compile

import (
	"Builder/artifact"
	"Builder/logger"
	"Builder/yaml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Go creates exe from file passed in as arg
func Go(filePath string) {

	//Set default project type env for builder.yaml creation
  projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "go")
	}

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
		os.Setenv("BUILDER_DIR_PATH", path)
	}
	
	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	//find 'go file' to be built
	buildFile := strings.ToLower(os.Getenv("BUILDER_BUILD_FILE"))
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	//if no file defined by user, use default main.go
	if (buildFile == "") {
		buildFile = "main.go"
		os.Setenv("BUILDER_BUILD_FILE", buildFile)
	}

	//buildName = buildfile (get rid of ".go") + Unix timestamp
	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		cmd = exec.Command(buildCmd)
	} else if (buildTool == "go") {
		fmt.Println(buildTool)
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath       // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath       // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "go build "+buildFile)
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Go project failed to compile.")
		log.Fatal(err)
	}

	yaml.CreateBuilderYaml(fullPath)
	
	//rename artifact by adding Unix timestamp
	_, extName := artifact.ExtExistsFunction(fullPath, ".exe")
	artifactName := artifact.NameArtifact(fullPath, extName)

	//send artifact to user specified path
	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", fullPath+"/"+artifactName, artifactPath).Run()
	}

	logger.InfoLogger.Println("Go project compiled successfully.")
}

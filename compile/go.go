// takes in code as arg from go
//run go build on code given

package compile

import (
	"Builder/artifact"
	"Builder/logger"
	"Builder/utils"
	"Builder/yaml"
	"bytes"
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

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	//find 'go file' to be built
	buildFile := strings.ToLower(os.Getenv("BUILDER_BUILD_FILE"))
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	//if no file defined by user, use default main.go
	if buildFile == "" {
		buildFile = "main.go"
		os.Setenv("BUILDER_BUILD_FILE", buildFile)
	}

	//buildName = buildfile (get rid of ".go") + Unix timestamp
	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		buildCmdArray := strings.Fields(buildCmd)
		cmd = exec.Command(buildCmdArray[0], buildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else if buildTool == "go" {
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "go build "+buildFile)
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		logger.ErrorLogger.Println("Go project failed to compile.")
		fmt.Println("out:", outb.String(), "err:", errb.String())
		log.Fatal(err)
	}
	yaml.CreateBuilderYaml(fullPath)

	packageGoArtifact(fullPath)

	logger.InfoLogger.Println("Go project compiled successfully.")
}

func packageGoArtifact(fullPath string) {
	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".exe")
	//copy artifact, then remove artifact in workspace
	exec.Command("cp", "-a", fullPath+"/"+extName, artifactDir).Run()
	exec.Command("rm", fullPath+"/"+extName).Run()

	//create metadata, then copy contents to zip dir
	utils.Metadata(artifactDir)
	artifact.ZipArtifactDir()

	//copy zip into open artifactDir, delete zip in workspace (keeps entire artifact contained)
	exec.Command("cp", "-a", artifactDir+".zip", artifactDir).Run()
	exec.Command("rm", artifactDir+".zip").Run()

	// artifactName := artifact.NameArtifact(fullPath, extName)

	// send artifact to user specified path
	artifactStamp := os.Getenv("BUILDER_ARTIFACT_STAMP")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if outputPath != "" {
		exec.Command("cp", "-a", artifactDir+"/"+artifactStamp+".zip", outputPath).Run()
	}
}

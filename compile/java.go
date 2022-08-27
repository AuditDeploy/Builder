package compile

import (
	"builder/artifact"
	"builder/logger"
	"builder/utils"
	"builder/yaml"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Java does ...
func Java(filePath string) {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("builder_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("builder_PROJECT_TYPE", "java")
	}

	//define dir path for command to run in
	var fullPath string
	configPath := os.Getenv("builder_DIR_PATH")
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
		os.Setenv("builder_DIR_PATH", path)

	}

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("builder_BUILD_TOOL"))
	buildCmd := os.Getenv("builder_BUILD_COMMAND")

	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		buildCmdArray := strings.Fields(buildCmd)
		cmd = exec.Command(buildCmdArray[0], buildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else if buildTool == "maven" || buildTool == "mvn" {
		fmt.Println(buildTool)
		cmd = exec.Command("mvn", "clean", "install")
		cmd.Dir = fullPath // or whatever directory it's in
	} else if buildTool == "gradle" {
		// gradle, etc.
	} else {
		//default
		cmd = exec.Command("mvn", "clean", "install")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("builder_BUILD_TOOL", "maven")
		os.Setenv("builder_BUILD_COMMAND", "mvn clean install")
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		logger.ErrorLogger.Println("Java project failed to compile.")
		fmt.Println("out:", outb.String(), "err:", errb.String())
		log.Fatal(err)
	}

	//creates default builder.yaml if it doesn't exist
	yaml.CreatebuilderYaml(fullPath)

	packageJavaArtifact(fullPath + "/target")

	logger.InfoLogger.Println("Java project compiled successfully.")
}
func packageJavaArtifact(fullPath string) {
	artifact.ArtifactDir()
	artifactDir := os.Getenv("builder_ARTIFACT_DIR")
	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".jar")
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
	artifactStamp := os.Getenv("builder_ARTIFACT_STAMP")
	outputPath := os.Getenv("builder_OUTPUT_PATH")
	if outputPath != "" {
		exec.Command("cp", "-a", artifactDir+"/"+artifactStamp+".zip", outputPath).Run()
	}
}

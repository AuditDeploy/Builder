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
	"path/filepath"
	"strings"
)

func CSharp(filePath string) {
	fmt.Println("C# filePath: " + filePath)
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("builder_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("builder_PROJECT_TYPE", "c#")
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

	//install dependencies/build,
	// if yaml build type exists install accordingly, if buildCmd exists,
	buildTool := strings.ToLower(os.Getenv("builder_BUILD_TOOL"))
	buildCmd := os.Getenv("builder_BUILD_COMMAND")

	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		buildCmdArray := strings.Fields(buildCmd)
		cmd = exec.Command(buildCmdArray[0], buildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else if buildTool == "dotnet" {
		cmd = exec.Command("dotnet", "build", fullPath)
		cmd.Dir = fullPath // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("dotnet", "build", fullPath)
		// cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("builder_BUILD_TOOL", "dotnet")
		os.Setenv("builder_BUILD_COMMAND", "dotnet build "+fullPath)
		os.Setenv("builder_BUILD_FILE", fullPath[strings.LastIndex(fullPath, "/")+1:])
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		logger.ErrorLogger.Println("C# project failed to compile.")
		fmt.Println("out:", outb.String(), "err:", errb.String())
		log.Fatal(err)
	}

	fullPath = fullPath[:strings.LastIndex(fullPath, "/")+1]

	yaml.CreatebuilderYaml(fullPath)

	packageCSharpArtifact(fullPath)

	logger.InfoLogger.Println("C# project compiled successfully.")
}

func packageCSharpArtifact(fullPath string) {
	artifact.ArtifactDir()
	artifactDir := os.Getenv("builder_ARTIFACT_DIR")
	//find artifact by extension
	artifactsArray, _ := WalkMatch(fullPath, "*.dll")

	//copy artifact, then remove artifact in workspace
	exec.Command("cp", "-a", artifactsArray[0], artifactDir).Run()
	exec.Command("rm", artifactsArray[0]).Run()

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

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

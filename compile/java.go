package compile

import (
	"Builder/artifact"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var locallogger *zap.Logger

// Java does ...
func Java(filePath string) {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "java")
	}

	//Set up local logger
        localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
        locallogger = log.NewLogger("logs", localPath)

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
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

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
		os.Setenv("BUILDER_BUILD_TOOL", "maven")
		os.Setenv("BUILDER_BUILD_COMMAND", "mvn clean install")
	}

	//run cmd, check for err, log cmd
	log.Info("run command", cmd)
	out, err := cmd.Output()

        //Log output to local log
        stdOut := string(out[:])
        if stdOut != "" {
                locallogger.Info(stdOut)
        }

        if err != nil {
                var outb, errb bytes.Buffer
                locallogger.Error(errb.String())
                log.Fatal("JAVA project failed to build", err)
        }

	//creates default builder.yaml if it doesn't exist
	yaml.CreateBuilderYaml(fullPath)

	packageJavaArtifact(fullPath + "/target")

	log.Info("Java project compiled successfully.")
}
func packageJavaArtifact(fullPath string) {
	archiveExt := ""

	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".jar")
	//copy artifact, then remove artifact in workspace
	exec.Command("cp", "-a", fullPath+"/"+extName, artifactDir).Run()
	exec.Command("rm", fullPath+"/"+extName).Run()

	//create metadata, then copy contents to zip dir
	utils.Metadata(artifactDir)

	if os.Getenv("ARTIFACT_ZIP_ENABLED") == "true" {
		//zip artifact
		artifact.ZipArtifactDir()

		//copy zip into open artifactDir, delete zip in workspace (keeps entire artifact contained)
		exec.Command("cp", "-a", artifactDir+archiveExt, artifactDir).Run()
		exec.Command("rm", artifactDir+archiveExt).Run()

		// artifactName := artifact.NameArtifact(fullPath, extName)

		// send artifact to user specified path or send to parent directory
		artifactStamp := os.Getenv("BUILDER_ARTIFACT_STAMP")
		outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
		if outputPath != "" {
			exec.Command("cp", "-a", artifactDir+"/"+artifactStamp+archiveExt, outputPath).Run()
		} else {
			exec.Command("cp", "-a", artifactDir+"/"+artifactStamp+archiveExt, os.Getenv("BUILDER_PARENT_DIR")).Run()
		}

		//remove artifact directory
		exec.Command("rm", "-r", artifactDir).Run()
	}
}

package compile

import (
	"Builder/artifact"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
)

// Java does ...
func Java(filePath string) {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "java")
	}

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger, closeLocalLogger = log.NewLogger("logs", localPath)

	//define dir path for command to run in
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included already, else add curren dir + local path
	if os.Getenv("BUILDER_COMMAND") == "true" {
		// ex: C:/Users/Name/Projects/helloworld_19293/workspace/dir
		fullPath = filePath
	} else if configPath != "" {
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
		os.Setenv("BUILDER_BUILD_COMMAND", "mvn clean install")
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
	spinner.LogMessage("running command: "+cmd.String(), "info")

	stdout, pipeErr := cmd.StdoutPipe()
	if pipeErr != nil {
		spinner.LogMessage(pipeErr.Error(), "fatal")
	}

	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	scanner := bufio.NewScanner(stdout)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			spinner.Spinner.Stop()
			locallogger.Info(line)
			spinner.Spinner.Start()
		}

		// We're all done, unblock the channel
		done <- struct{}{}

	}()

	os.Setenv("BUILD_START_TIME", time.Now().Format(time.RFC850))

	if err := cmd.Start(); err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	// Wait for all output to be processed
	<-done

	// Wait for cmd to finish
	if err := cmd.Wait(); err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	os.Setenv("BUILD_END_TIME", time.Now().Format(time.RFC850))

	// Close log file
	closeLocalLogger()

	// Update parent dir name to include start time and send back new full path
	fullPath = directory.UpdateParentDirName(fullPath)

	//creates default builder.yaml if it doesn't exist
	if os.Args[1] == "init" || os.Args[1] == "config" {
		repoPath := "./" + strings.TrimSuffix(utils.GetName(), ".git")

		if configPath != "" {
			repoPath = configPath + "/" + strings.TrimSuffix(utils.GetName(), ".git")
		}

		yaml.CreateBuilderYaml(repoPath)
	} else {
		yaml.CreateBuilderYaml(fullPath)
	}

	packageJavaArtifact(fullPath + "/target")

	spinner.LogMessage("Java project compiled successfully.", "info")
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
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")

	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".jar")
	os.Setenv("BUILDER_ARTIFACT_NAMES", extName)

	//copy artifact, then remove artifact in workspace
	err := cp.Copy(fullPath+"/"+extName, artifactDir+"/"+extName)
	if err != nil {
		spinner.LogMessage(err.Error(), "warn")
	}

	// If outputpath provided also cp artifacts to that location
	if outputPath != "" {
		// Check if outputPath exists.  If not, create it
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			if err := os.Mkdir(outputPath, 0755); err != nil {
				spinner.LogMessage("Could not create output path", "fatal")
			}
		}

		err := cp.Copy(fullPath+"/"+extName, outputPath+"/"+extName)
		if err != nil {
			spinner.LogMessage(err.Error(), "warn")
		}

		spinner.LogMessage("Artifact(s) copied to output path provided", "info")
	}

	errRemove := os.Remove(fullPath + "/" + extName)
	if errRemove != nil {
		spinner.LogMessage(errRemove.Error(), "warn")
	}

	//create metadata, then copy contents to zip dir
	utils.Metadata(artifactDir)

	if os.Getenv("ARTIFACT_ZIP_ENABLED") == "true" {
		//zip artifact
		artifact.ZipArtifactDir()

		//remove uncompressed artifact
		err := os.Remove(artifactDir + "/" + extName)
		if err != nil {
			spinner.LogMessage(err.Error(), "warn")
		}

		// send artifact to user specified path or send to artifact directory
		outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
		if outputPath != "" {
			err := cp.Copy(artifactDir+archiveExt, outputPath+"/"+filepath.Base(artifactDir)+archiveExt)
			if err != nil {
				spinner.LogMessage(err.Error(), "warn")
			}
		} else {
			err := cp.Copy(artifactDir+archiveExt, artifactDir+"/"+filepath.Base(artifactDir)+archiveExt)
			if err != nil {
				spinner.LogMessage(err.Error(), "warn")
			}
		}

		errRemove := os.Remove(artifactDir + archiveExt)
		if errRemove != nil {
			spinner.LogMessage(errRemove.Error(), "warn")
		}
	}
}

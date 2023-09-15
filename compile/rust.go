// takes in code as arg from rust
//run rust build on code given

package compile

import (
	"Builder/artifact"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Rust creates exe from file passed in as arg
func Rust(filePath string) {

	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "rust")
	}

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger, closeLocalLogger = log.NewLogger("logs", localPath)

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
	//find 'rs file' to be built
	buildFile := os.Getenv("BUILDER_BUILD_FILE")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	//if no file defined by user, use default Cargo.toml
	if buildFile == "" {
		buildFile = "Cargo.toml"
		os.Setenv("BUILDER_BUILD_FILE", buildFile)
	}

	//buildName = buildfile (get rid of ".rs") + Unix timestamp
	var cmd *exec.Cmd
	if buildCmd != "" {
		//user specified cmd
		buildCmdArray := strings.Fields(buildCmd)
		cmd = exec.Command(buildCmdArray[0], buildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else if buildTool == "rust" {
		cmd = exec.Command("cargo", "build", "-r")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "cargo build -r")
	} else {
		cmd = exec.Command("cargo", "build", "-r")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "cargo build -r")
		os.Setenv("BUILDER_BUILD_TOOL", "rust")
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

	// Update parent dir name to include start time
	fullPath = directory.UpdateParentDirName(fullPath)

	yaml.CreateBuilderYaml(fullPath)

	packageRustArtifact(fullPath)

	spinner.LogMessage("Go project built successfully.", "info")
}

func packageRustArtifact(fullPath string) {
	archiveExt := ""
	artifactExt := ""
	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
		artifactExt = ".exe"
	} else {
		archiveExt = ".tar.gz"
		artifactExt = ""
	}

	extName := ""

	tomlfile, _ := os.Open(fullPath+"/" + os.Getenv("BUILDER_BUILD_FILE"))
	scanner := bufio.NewScanner(tomlfile)
	for scanner.Scan() {
        line:= scanner.Text()
		if strings.HasPrefix(line, "name = ") == true {
			extName = strings.Trim(line[7:]+artifactExt, "\"")
		 } 
    }
	defer tomlfile.Close()

	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")

	//find artifact by extension
	os.Setenv("BUILDER_ARTIFACT_NAMES", extName)
	//copy artifact, then remove artifact in workspace
	exec.Command("cp", "-a", fullPath+"/target/release/"+extName, artifactDir).Run()

	// If outputpath provided also cp artifacts to that location
	if outputPath != "" {
		// Check if outputPath exists.  If not, create it
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			if err := os.Mkdir(outputPath, 0755); err != nil {
				spinner.LogMessage("Could not create output path", "fatal")
			}
		}

		exec.Command("cp", "-a", fullPath+"/target/release/"+extName, outputPath).Run()

		spinner.LogMessage("Artifact(s) copied to output path provided", "info")
	}

	exec.Command("rm", fullPath+"/target/release/"+extName).Run()

	//create metadata
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

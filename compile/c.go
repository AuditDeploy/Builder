package compile

import (
	"Builder/artifact"
	"Builder/directory"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()
var closeLocalLogger func()

// C/C++ does ...
func C(filePath string) {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "c")
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
	//find 'Makefile' to be built
	buildFile := strings.ToLower(os.Getenv("BUILDER_BUILD_FILE"))
	configCmd := os.Getenv("BUILDER_CONFIG_COMMAND")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

	var cmd *exec.Cmd

	if configCmd != "" {
		//user specified cmd
		configCmdArray := strings.Fields(configCmd)
		cmd = exec.Command(configCmdArray[0], configCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("./configure")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_CONFIG_COMMAND", "./configure")
	}

	//run config cmd, check for err, log config cmd
	BuilderLog.Infof("run command", cmd)
	err := cmd.Run()
	if err != nil {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		fmt.Println("out:", outb.String(), "err:", errb.String())
		BuilderLog.Fatalf("C/C++ failed to compile", err)
	}

	if buildCmd != "" {
		//user specified cmd
		buildCmdArray := strings.Fields(buildCmd)
		cmd = exec.Command(buildCmdArray[0], buildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in
	} else if strings.Contains(buildTool, "Make") && buildFile != "" {
		cmd = exec.Command("make -f", buildFile)
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "make -f "+buildFile)
	} else {
		//default
		cmd = exec.Command("make")
		cmd.Dir = fullPath   // or whatever directory it's in
		if buildTool == "" { // If buildTool hasn't been set yet, set it
			os.Setenv("BUILDER_BUILD_TOOL", "Make")
		}
		os.Setenv("BUILDER_BUILD_COMMAND", "make")
	}

	//run cmd, check for err, log cmd
	BuilderLog.Infof("running command: ", os.Getenv("BUILDER_BUILD_COMMAND"))

	stdout, pipeErr := cmd.StdoutPipe()
	if pipeErr != nil {
		BuilderLog.Fatal(pipeErr.Error())
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
			locallogger.Info(line)
		}

		// We're all done, unblock the channel
		done <- struct{}{}

	}()

	os.Setenv("BUILD_START_TIME", time.Now().Format(time.RFC850))

	if err := cmd.Start(); err != nil {
		BuilderLog.Fatal(err.Error())
	}

	// Wait for all output to be processed
	<-done

	// Wait for cmd to finish
	if err := cmd.Wait(); err != nil {
		BuilderLog.Fatal(err.Error())
	}

	os.Setenv("BUILD_END_TIME", time.Now().Format(time.RFC850))

	// Close log file
	closeLocalLogger()

	// Update parent dir name to include start time and send back new full path
	fullPath = directory.UpdateParentDirName(fullPath)

	//creates default builder.yaml if it doesn't exist
	yaml.CreateBuilderYaml(fullPath)

	packageCArtifact(fullPath + "/build")

	BuilderLog.Info("C/C++ project compiled successfully.")
}

func packageCArtifact(fullPath string) {
	archiveExt := ""

	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	artifactList := os.Getenv("BUILDER_ARTIFACT_LIST")

	// If we were given an artifacts list, handle it
	if artifactList != "" {
		artifactArray := strings.Fields(artifactList)

		//copy artifact(s), then remove artifact(s) from workspace
		for _, artifact := range artifactArray {
			exec.Command("cp", "-a", artifact, artifactDir).Run()
			exec.Command("rm", artifact).Run()
		}

	} else {
		var artifactExt string
		buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
		//Determine artifact extension
		switch buildTool {
		case "Make-rpm":
			artifactExt = "*.rpm"
		case "Make-deb":
			artifactExt = "*.deb"
		case "Make-tar":
			artifactExt = "*.tar.gz"
		case "Make-lib":
			archiveExt = "*.lib"
		case "Make-dll":
			archiveExt = "*.dll"
		default:
			artifactExt = "*.exe"
		}

		//find artifact(s) by extension
		// WalkMatch function defined in compile/c#.go
		artifactArray, _ := WalkMatch(fullPath, artifactExt)

		//copy artifact(s), then remove artifact(s) from workspace
		for i := 0; i < len(artifactArray); i++ {
			exec.Command("cp", "-a", artifactArray[i], artifactDir).Run()
			exec.Command("rm", artifactArray[i]).Run()
		}
	}

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

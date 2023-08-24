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
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

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
	preBuildCmd := os.Getenv("BUILDER_PREBUILD_COMMAND")
	configCmd := os.Getenv("BUILDER_CONFIG_COMMAND")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

	var cmd *exec.Cmd

	// If a pre-build command is provided execute it
	if preBuildCmd != "" {
		//user specified cmd
		preBuildCmdArray := strings.Fields(preBuildCmd)
		cmd = exec.Command(preBuildCmdArray[0], preBuildCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in

		//run pre-build cmd, check for err, log pre-build cmd
		spinner.LogMessage("running command: "+cmd.String(), "info")

		preBuildStdout, pipeErr := cmd.StdoutPipe()
		if pipeErr != nil {
			spinner.LogMessage(pipeErr.Error(), "fatal")
		}

		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		preBuildDone := make(chan struct{})

		preBuildScanner := bufio.NewScanner(preBuildStdout)

		// Use the scanner to scan the output line by line and log it
		// It's running in a goroutine so that it doesn't block
		go func() {
			// Read line by line and process it
			for preBuildScanner.Scan() {
				line := preBuildScanner.Text()
				// Have to stop spinner or it will get printed with log to console
				spinner.Spinner.Stop()
				locallogger.Info(line)
				spinner.Spinner.Start()
			}

			// We're all done, unblock the channel
			preBuildDone <- struct{}{}

		}()

		if err := cmd.Start(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Wait for all output to be processed
		<-preBuildDone

		// Wait for cmd to finish
		if err := cmd.Wait(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}
	}

	// If a configure command is provided execute it
	if configCmd != "" {
		//user specified cmd
		configCmdArray := strings.Fields(configCmd)
		cmd = exec.Command(configCmdArray[0], configCmdArray[1:]...)
		cmd.Dir = fullPath // or whatever directory it's in

		//run config cmd, check for err, log config cmd
		spinner.LogMessage("running command: "+cmd.String(), "info")

		configStdout, pipeErr := cmd.StdoutPipe()
		if pipeErr != nil {
			spinner.LogMessage(pipeErr.Error(), "fatal")
		}

		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		configDone := make(chan struct{})

		configScanner := bufio.NewScanner(configStdout)

		// Use the scanner to scan the output line by line and log it
		// It's running in a goroutine so that it doesn't block
		go func() {
			// Read line by line and process it
			for configScanner.Scan() {
				line := configScanner.Text()
				// Have to stop spinner or it will get printed with log to console
				spinner.Spinner.Stop()
				locallogger.Info(line)
				spinner.Spinner.Start()
			}

			// We're all done, unblock the channel
			configDone <- struct{}{}

		}()

		if err := cmd.Start(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Wait for all output to be processed
		<-configDone

		// Wait for cmd to finish
		if err := cmd.Wait(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}
	}

	// Build command
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
			// Have to stop spinner or it will get printed with log to console
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
	yaml.CreateBuilderYaml(fullPath)

	packageCArtifact(fullPath)

	spinner.LogMessage("C/C++ project compiled successfully.", "info")
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
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")

	// If we were given an artifacts list, handle it
	if artifactList != "" {
		artifactArray := strings.Split(artifactList, ",")
		os.Setenv("BUILDER_ARTIFACT_NAMES", artifactList)

		//copy artifact(s), then remove artifact(s) from workspace
		for _, artifact := range artifactArray {
			exec.Command("cp", fullPath+"/"+artifact, artifactDir).Run()

			// If outputpath provided also cp artifacts to that location
			if outputPath != "" {
				// Check if outputPath exists.  If not, create it
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					if err := os.Mkdir(outputPath, 0755); err != nil {
						spinner.LogMessage("Could not create output path", "fatal")
					}
				}

				exec.Command("cp", fullPath+"/"+artifact, outputPath).Run()

				spinner.LogMessage("Artifact(s) copied to output path provided", "info")
			}

			exec.Command("rm", fullPath+"/"+artifact).Run()
		}

	} else {
		var artifactExt string
		buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
		//Determine artifact extension
		switch buildTool {
		case "make-rpm":
			artifactExt = "*.rpm"
		case "make-deb":
			artifactExt = "*.deb"
		case "make-tar":
			artifactExt = "*.tar.gz"
		case "make-lib":
			archiveExt = "*.lib"
		case "make-dll":
			archiveExt = "*.dll"
		default:
			artifactExt = "*.exe"
		}

		//find artifact(s) by extension
		// WalkMatch function defined in compile/c#.go
		artifactArray, _ := WalkMatch(fullPath, artifactExt)
		if len(artifactArray) == 0 {
			spinner.LogMessage("Could not find artifact(s).  Please specify the name(s) in the artifactlist of the builder.yaml", "fatal")
		}

		var artifactNames []string

		//copy artifact(s), then remove artifact(s) from workspace
		for i := 0; i < len(artifactArray); i++ {
			artifactNames = append(artifactNames, filepath.Base(artifactArray[i]))
			exec.Command("cp", artifactArray[i], artifactDir).Run()

			// If outputpath provided also cp artifacts to that location
			if outputPath != "" {
				// Check if outputPath exists.  If not, create it
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					if err := os.Mkdir(outputPath, 0755); err != nil {
						spinner.LogMessage("Could not create output path", "fatal")
					}
				}

				exec.Command("cp", artifactArray[i], outputPath).Run()

				spinner.LogMessage("Artifact(s) copied to output path provided", "info")
			}

			exec.Command("rm", artifactArray[i]).Run()
		}

		os.Setenv("BUILDER_ARTIFACT_NAMES", strings.Join([]string(artifactNames), ","))
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

package compile

import (
	"Builder/artifact"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
)

// Ruby creates zip from files passed in as arg
func Ruby() {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "ruby")
	}

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger, closeLocalLogger = log.NewLogger("logs", localPath)

	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	tempWorkspace := workspaceDir + "/temp/"
	//make temp dir
	os.Mkdir(tempWorkspace, 0755)

	//add hidden dir contents to temp dir, install dependencies
	err := cp.Copy(hiddenDir+"/.", tempWorkspace)
	if err != nil {
		spinner.LogMessage(err.Error(), "warn")
	}

	//define dir path for command to run
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included in tempWorkspace, else add the local path
	if os.Getenv("BUILDER_COMMAND") == "true" {
		// ex: C:/Users/Name/Projects/helloworld_19293/workspace/dir
		fullPath = tempWorkspace
	} else if configPath != "" {
		fullPath = tempWorkspace
	} else {
		path, _ := os.Getwd()
		//combine local path to newly created tempWorkspace, gets rid of "." in path name
		fullPath = path + tempWorkspace[strings.Index(tempWorkspace, ".")+1:]
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
	} else if buildTool == "Bundler" {
		fmt.Println(buildTool)
		cmd = exec.Command("bundle", "install", "--path", "vendor/bundle")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "bundle install --path vendor/bundle")
	} else {
		//default
		cmd = exec.Command("bundle", "install", "--path", "vendor/bundle")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_TOOL", "bundler")
		os.Setenv("BUILDER_BUILD_COMMAND", "bundle install --path vendor/bundle")
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

	// Update vars because of parent dir name change
	hiddenDir = os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir = os.Getenv("BUILDER_WORKSPACE_DIR")
	tempWorkspace = workspaceDir + "/temp/"

	if os.Args[1] == "init" || os.Args[1] == "config" {
		repoPath := "./" + strings.TrimSuffix(utils.GetName(), ".git")

		if configPath != "" {
			repoPath = configPath + "/" + strings.TrimSuffix(utils.GetName(), ".git")
		}

		yaml.CreateBuilderYaml(repoPath)
	} else {
		yaml.CreateBuilderYaml(fullPath)
	}

	//CreateZip artifact dir with timestamp
	parsedStartTime, _ := time.Parse(time.RFC850, os.Getenv("BUILD_START_TIME"))
	timeBuildStarted := parsedStartTime.Unix()

	outFile, err := os.Create(workspaceDir + "/artifact_" + strconv.FormatInt(timeBuildStarted, 10) + ".zip")
	if err != nil {
		spinner.LogMessage("Ruby failed to get artifact: "+err.Error(), "fatal")
	}

	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addRubyFiles(w, tempWorkspace, "")

	err = w.Close()
	if err != nil {
		spinner.LogMessage("Ruby project failed to compile: "+err.Error(), "fatal")
	}
	packageRubyArtifact(fullPath)

	// artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	// if artifactPath != "" {
	// 	exec.Command("cp", "-a", workspaceDir+"/temp.zip", artifactPath).Run()
	// }
	spinner.LogMessage("Ruby project compiled successfully.", "info")
}

func packageRubyArtifact(fullPath string) {
	archiveExt := ""

	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")

	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(workspaceDir, ".zip")
	os.Setenv("BUILDER_ARTIFACT_NAMES", extName)

	//copy artifact, then remove artifact in workspace
	err := cp.Copy(workspaceDir+"/"+extName, artifactDir+"/"+extName)
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

		err := cp.Copy(workspaceDir+"/"+extName, outputPath+"/"+extName)
		if err != nil {
			spinner.LogMessage(err.Error(), "warn")
		}

		spinner.LogMessage("Artifact(s) copied to output path provided", "info")
	}

	errRemove := os.Remove(workspaceDir + "/" + extName)
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

// recursively add files
func addRubyFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			addRubyFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

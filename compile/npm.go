package compile

import (
	"Builder/artifact"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Npm creates zip from files passed in as arg
func Npm() {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "node")
	}

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger = log.NewLogger("logs", localPath)

	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	tempWorkspace := workspaceDir + "/temp/"
	//make temp dir
	os.Mkdir(tempWorkspace, 0755)

	//add hidden dir contents to temp dir, install dependencies
	exec.Command("cp", "-a", hiddenDir+"/.", tempWorkspace).Run()

	//define dir path for command to run
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included in tempWorkspace, else add the local path
	if configPath != "" {
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
	} else if buildTool == "npm" {
		cmd = exec.Command("npm", "install")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_COMMAND", "npm install")
	} else {
		//default
		cmd = exec.Command("npm", "install")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_TOOL", "npm")
		os.Setenv("BUILDER_BUILD_COMMAND", "npm install")
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

	if err := cmd.Start(); err != nil {
		BuilderLog.Fatal(err.Error())
	}

	// Wait for all output to be processed
	<-done

	// Wait for cmd to finish
	if err := cmd.Wait(); err != nil {
		BuilderLog.Fatal(err.Error())
	}

	yaml.CreateBuilderYaml(fullPath)

	//sets path for metadata, and addFiles (covers when wrkspace dir env doesn't exist)
	var addPath string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		path, _ := os.Getwd()
		addPath = path + "/"
	} else {
		addPath = tempWorkspace
	}

	utils.Metadata(addPath)

	//sets path for zip creation
	var dirPath string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		path, _ := os.Getwd()
		dirPath = strings.Replace(path, "\\temp", "", 1)
	} else {
		dirPath = workspaceDir
	}

	// CreateZip artifact dir with timestamp
	currentTime := time.Now().Unix()

	outFile, err := os.Create(dirPath + "/artifact_" + strconv.FormatInt(currentTime, 10) + ".zip")
	if err != nil {
		BuilderLog.Fatalf("node-npm failed to get artifact", err)
	}

	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addNpmFiles(w, addPath, "")

	err = w.Close()
	if err != nil {
		BuilderLog.Fatalf("node-npm project failed to compile", err)
	}

	packageNpmArtifact(fullPath)
	// artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	// fmt.Print(artifactPath)
	// if artifactPath != "" {
	// 	artifactZip := os.Getenv("BUILDER_ARTIFACT_STAMP")
	// 	fmt.Print(artifactZip)
	// 	exec.Command("cp", "-a", artifactZip+".zip", artifactPath).Run()
	// }
	BuilderLog.Info("node-npm project compiled successfully")
}

func packageNpmArtifact(fullPath string) {
	archiveExt := ""

	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".exe")
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

// recursively add files
func addNpmFiles(w *zip.Writer, basePath, baseInZip string) {
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
			addNpmFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

package compile

import (
	"Builder/artifact"
	"Builder/logger"
	"Builder/utils"
	"Builder/yaml"
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	//tempWorkspace := workspaceDir + "/temp/"
	tempWorkspace := filepath.Join(workspaceDir, "temp")
	//make temp dir
	os.Mkdir(tempWorkspace, 0755)

	// get current working directory
	path, _ := os.Getwd()
	//add hidden dir contents to temp dir, install dependencies
	//exec.Command("cp", "-a", path+hiddenDir+"/.", tempWorkspace).Run()
	if runtime.GOOS == "windows" {
		logger.InfoLogger.Println("this is windows")
		exec.Command("cmd", "/C", "xcopy", filepath.Join(path, hiddenDir), filepath.Join(path, tempWorkspace), "/e").Run()
	} else {
		exec.Command("cp", "-a", filepath.Join(path, hiddenDir), filepath.Join(path, tempWorkspace)).Run()
	}

	//define dir path for command to run
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included in tempWorkspace, else add the local path
	if configPath != "" {
		fullPath = tempWorkspace
	} else {
		path, _ := os.Getwd()
		//combine local path to newly created tempWorkspace, gets rid of "." in path name
		//fullPath = path + tempWorkspace[strings.Index(tempWorkspace, ".")+1:]
		fullPath = filepath.Join(path, tempWorkspace)

		//testing changing default path to hidden folder
		//fullPath = hiddenDir + "\\.hidden"
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
	} else {
		//default
		if !(runtime.GOOS == "windows") {
			fullPath = filepath.Join(fullPath, ".hidden")
		}
		cmd = exec.Command("npm", "install")
		cmd.Dir = fullPath // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_TOOL", "npm")
		os.Setenv("BUILDER_BUILD_COMMAND", "npm install")
	}

	//run cmd, check for err, log cmd
	err := cmd.Run()
	if err != nil {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		logger.ErrorLogger.Println("Node project failed to compile.")
		fmt.Println("out:", outb.String(), "err:", errb.String())
		log.Fatal(err)
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
		log.Fatal(err)
	}

	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addNpmFiles(w, addPath, "")

	err = w.Close()
	if err != nil {
		logger.ErrorLogger.Println("Npm project failed to compile.")
		log.Fatal(err)
	}

	packageNpmArtifact(fullPath)
	// artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	// fmt.Print(artifactPath)
	// if artifactPath != "" {
	// 	artifactZip := os.Getenv("BUILDER_ARTIFACT_STAMP")
	// 	fmt.Print(artifactZip)
	// 	exec.Command("cp", "-a", artifactZip+".zip", artifactPath).Run()
	// }
	logger.InfoLogger.Println("Npm project compiled successfully.")
}

func packageNpmArtifact(fullPath string) {
	artifact.ArtifactDir()
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	//find artifact by extension
	_, extName := artifact.ExtExistsFunction(fullPath, ".exe")
	//copy artifact, then remove artifact in workspace
	currentDirectory, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "xcopy", filepath.Join(fullPath, extName), filepath.Join(currentDirectory, artifactDir), "/e").Run()
		exec.Command("cmd", "/C", "del", filepath.Join(fullPath, extName)).Run()
	} else {
		//exec.Command("cp", "-a", fullPath+"/"+extName, artifactDir).Run()
		//exec.Command("rm", fullPath+"/"+extName).Run()
		exec.Command("cp", "-a", filepath.Join(fullPath, extName), filepath.Join(currentDirectory, artifactDir)).Run()
		exec.Command("rm", filepath.Join(fullPath, extName)).Run()
	}

	//create metadata, then copy contents to zip dir
	utils.Metadata(artifactDir)
	artifact.ZipArtifactDir()

	//copy zip into open artifactDir, delete zip in workspace (keeps entire artifact contained)
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "xcopy", artifactDir+".zip", artifactDir, "/e").Run()
		exec.Command("cmd", "/C", "del", artifactDir+".zip").Run()
	} else {
		exec.Command("cp", "-a", artifactDir+".zip", artifactDir).Run()
		exec.Command("rm", artifactDir+".zip").Run()
	}

	// artifactName := artifact.NameArtifact(fullPath, extName)

	// send artifact to user specified path
	artifactStamp := os.Getenv("BUILDER_ARTIFACT_STAMP")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if outputPath != "" {
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "xcopy", filepath.Join(artifactDir, artifactStamp+".zip"), outputPath, "/e").Run()
		} else {
			exec.Command("cp", "-a", artifactDir+"/"+artifactStamp+".zip", outputPath).Run()
		}
	}
}

// recursively add files
func addNpmFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	//get current working directory
	currentDir, _ := os.Getwd()

	for _, file := range files {
		if !file.IsDir() {
			//dat, err := ioutil.ReadFile(basePath + file.Name())
			dat, err := ioutil.ReadFile(filepath.Join(currentDir, basePath, file.Name()))
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
			//newBase := basePath + file.Name() + "/"
			newBase := filepath.Join(basePath, file.Name())
			//addNpmFiles(w, newBase, baseInZip+file.Name()+"/")
			addNpmFiles(w, newBase, baseInZip+file.Name())
		}
	}
}

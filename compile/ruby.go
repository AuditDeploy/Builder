package compile

import (
	"Builder/logger"
	"Builder/utils"
	"Builder/yaml"
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//Ruby creates zip from files passed in as arg
func Ruby() {
	//Set default project type env for builder.yaml creation
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	if projectType == "" {
		os.Setenv("BUILDER_PROJECT_TYPE", "ruby")
	}

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
	if (configPath != "") {
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
	} else if (buildTool == "Bundler") {
		fmt.Println(buildTool)
		cmd = exec.Command("bundle", "install", "--path", "vendor/bundle")
    cmd.Dir = fullPath       // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("bundle", "install", "--path", "vendor/bundle") 
    cmd.Dir = fullPath       // or whatever directory it's in
		os.Setenv("BUILDER_BUILD_TOOL", "bundler")
		os.Setenv("BUILDER_BUILD_COMMAND", "bundle install --path vendor/bundle")
	}
	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Ruby project failed to compile.")
		fmt.Println(err)
		log.Fatal(err)
	}

	yaml.CreateBuilderYaml(fullPath)

	//sets path for metadata, and addFiles (covers when workspace dir env doesn't exist)
	var addPath string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		path, _ := os.Getwd()
		addPath = path+"/"
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
	
	//CreateZip artifact dir with timestamp
	currentTime := time.Now().Unix()

	outFile, err := os.Create(dirPath+"/artifact_"+strconv.FormatInt(currentTime, 10)+".zip")
	if err != nil {
		log.Fatal(err)
	}
	
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addRubyFiles(w, tempWorkspace, "")

	err = w.Close()
	if err != nil {
		logger.ErrorLogger.Println("Ruby project failed to compile.")
		 log.Fatal(err)
	}

	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", workspaceDir+"/temp.zip", artifactPath).Run()
	}
	logger.InfoLogger.Println("Ruby project compiled successfully.")
}

//recursively add files
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
					addRubyFiles(w, newBase, baseInZip  + file.Name() + "/")
			}
	}
}
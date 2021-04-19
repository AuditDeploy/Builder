package utils

import (
	"Builder/compile"
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

//ProjectType will derive the project type and execute its compiler
func ProjectType() {

	//check for user defined project type from builder.yaml to define string array files
	configType := strings.ToLower(os.Getenv("BUILDER_PROJECT_TYPE"))

	var files []string
	//projectType exists in builder.yaml
	if configType != "" {
		//check value of config type, return string array of language's build file/files
		files = ConfigDerive()
	} else {
		//default
		files = []string{"main.go", "package.json", "pom.xml", "gemfile.lock", "gemfile", "pipfile.lock"}
	}

	//look for those files inside hidden dir
	for _, file := range files {

		//recursively check for file in hidden dir, return path if found
		filePath := findPath(file)
		//double check it exists
		fileExists, err := fileExistsInDir(filePath)
		if err != nil {
			logger.ErrorLogger.Println("No Go, Npm, Ruby, Python or Java File Exists")
			log.Fatal(err)
		}
		//if file exists and filePath isn't empty, run conditional to find correct compiler
		if fileExists && filePath != "" && filePath != "./" {
			if file == "main.go" || configType == "go" {
				//executes go compiler
				finalPath := createFinalPath(filePath, file)
				CopyDir()
				logger.InfoLogger.Println("Go project detected")
				compile.Go(finalPath)
				return
			} else if file == "package.json" || configType == "node" || configType == "npm" {
				//executes node compiler
				logger.InfoLogger.Println("Npm project detected")
				compile.Npm()
				return
			} else if file == "pom.xml" || configType == "java" {
				//executes java compiler
				CopyDir()
				logger.InfoLogger.Println("Java project detected")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.Java(workspace)
				return
			} else if file == "gemfile.lock" || file == "gemfile" || configType == "ruby" {
				//executes ruby compiler
				logger.InfoLogger.Println("Ruby project detected")
				compile.Ruby()
				return
			} else if file == "pipfile.lock" || configType == "python" {
				//executes python compiler
				logger.InfoLogger.Println("Python project detected")
				compile.Python()
				return
			}
		}
	}
	deriveProjectByExtension()
}

//derive projects by Extensions
func deriveProjectByExtension() {

	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	extensions := []string{".csproj", ".sln"}

	var filePathsFoundInRepo []string

	//finds all the paths of files with ext .csproj and .sln and appends them to
	//filePathsFoundInRepo array
	for _, ext := range extensions {
		extExists, fileNameArray := extExistsFunction(parentDir+"/"+".hidden", ext)
		if extExists {
			for _, fileName := range fileNameArray {
				filePathsFoundInRepo = append(filePathsFoundInRepo, findPath(fileName))
			}
		}
	}

	//checks if there's more than more file to compile from and if it does, prompt user to select path
	if len(filePathsFoundInRepo) > 1 {
		pathToCompileFrom := strings.Replace(selectPathToCompileFrom(filePathsFoundInRepo), ".hidden", "workspace", 1)
		CopyDir()
		logger.InfoLogger.Println("C# project detected")
		compile.CSharp(pathToCompileFrom)
	} else {
		pathToCompileFrom := strings.Replace(filePathsFoundInRepo[0], ".hidden", "workspace", 1)
		CopyDir()
		logger.InfoLogger.Println("C# project detected")
		compile.CSharp(pathToCompileFrom)
	}
}

//takes in file, searches hiddenDir to find a match and returns path to file
func findPath(file string) string {
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	// if f.Name is == to file passed in "coolProject.go", filePath becomes the path that file exists in
	var filePath string
	err := filepath.Walk(hiddenDir, func(path string, f os.FileInfo, err error) error {
		if strings.EqualFold(f.Name(), file) {
			filePath = path
		}
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, filePath has fullPath included, if not, run it locally
	if configPath != "" {
		return filePath
	} else {
		return "./" + filePath
	}
}

//changes .hidden to workspace for langs that produce binary, get's rid of file name in path
func createFinalPath(path string, file string) string {
	workFilePath := strings.Replace(path, ".hidden", "workspace", 1)
	finalPath := strings.Replace(workFilePath, file, "", -1)

	return finalPath
}

//checks if file exists
func fileExistsInDir(path string) (bool, error) {
	//file exists
	_, err := os.Stat(path)
	//return true
	if err == nil {
		return true, nil
	}
	//return false
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func extExistsFunction(dirPath string, ext string) (bool, []string) {
	found := false

	d, err := os.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var fileNameArray []string

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ext {
				fileNameArray = append(fileNameArray, file.Name())
				found = true
			}
		}
	}

	return found, fileNameArray
}

func selectPathToCompileFrom(filePaths []string) string {
	prompt := promptui.Select{
		Label: "Select a Path To Compile From: ",
		Items: filePaths,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

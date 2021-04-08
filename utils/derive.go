package utils

import (
	"Builder/compile"
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//ProjectType will derive the project type and execute its compiler
func ProjectType() {

	//parentDir = the path of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")
	//check for user defined type from builder.yaml
	configType := os.Getenv("BUILDER_PROJECT_TYPE")

	var files []string
	//if type is anything besides "", define the files var instead of looking for it
	if (configType != "") {
		//check value of config type, return string array of languages build file/files
		files = ConfigDerive()
	}	else {
		//set files var to default
		files = []string{"main.go", "package.json", "pom.xml", "gemfile.lock", "pipfile.lock"}
	}

	//look for those files inside hidden dir
	for _, file := range files {

		filePath := parentDir + "/" + ".hidden" + "/" + file
		//check if the filepath exists
		fileExists, err := exists(filePath)
		if err != nil {
			logger.ErrorLogger.Println("No Go, Npm, Ruby, Python or Java File Exists")
			log.Fatal(err)
		}

		//if file exists run a swith statement
		if fileExists {
			switch file {
			case "main.go":
				//executes go compiler
				CopyDir()
				logger.InfoLogger.Println("Go project detected")
				compile.Go(filePath)
			case "package.json":
				//executes node compiler
				logger.InfoLogger.Println("Npm project detected")
				compile.Npm()
			case "pom.xml":
				//executes java compiler
				CopyDir()
				logger.InfoLogger.Println("Java project detected")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.Java(workspace)
			case "gemfile.lock":
				//executes ruby compiler
				logger.InfoLogger.Println("Ruby project detected")
				compile.Ruby()
			case "pipfile.lock":
				//executes python compiler
				logger.InfoLogger.Println("Python project detected")
				compile.Python()
			}
		}
	}
	deriveProjectByExtension()
}

//derive projects by Extensions
func deriveProjectByExtension() {
	//parentDir = the path of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	extensions := []string{".csproj", ".sln"}

	for _, ext := range extensions {
		extExists := extExists(parentDir+"/"+".hidden", ext)

		if extExists {
			switch ext {
			case ".csproj":
				CopyDir()
				logger.InfoLogger.Println("C# project detected Ext csproj")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.CSharp(workspace)
			case ".sln":
				CopyDir()
				logger.InfoLogger.Println("C# project detected Ext sln")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.CSharp(workspace)
			}
		}

	}

}

//checks if file exists
func exists(path string) (bool, error) {
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

func extExists(dirPath string, ext string) bool {
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

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ext {

				found = true
			}
		}
	}

	return found
}

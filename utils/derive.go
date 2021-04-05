package utils

import (
	"Builder/compile"
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//ProjectType will derive the poject type(go, node, java repo) and execute its compiler
func ProjectType() {

	//parentDir = the name of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	//languages we are currently compiling
	files := []string{"main.go", "package.json", "pom.xml", "gemfile.lock"}

	for _, file := range files {

		filePath := parentDir + "/" + ".hidden" + "/" + file

		//checking if the filepath exists
		fileExists, err := exists(filePath)

		if err != nil {
			logger.ErrorLogger.Println("No Go, Npm, Ruby or Java File Exists")
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
				break
			case "package.json":
				//executes node compiler
				logger.InfoLogger.Println("Npm project detected")
				compile.Npm()
				break
			case "pom.xml":
				//executes java compiler
				CopyDir()
				logger.InfoLogger.Println("Java project detected")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.Java(workspace)
				break

			}
		}

	}
	deriveProjectByExtension()
}

//derive projects by Extensions
func deriveProjectByExtension() {
	//parentDir = the name of the project
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
				break
			case ".sln":
				CopyDir()
				logger.InfoLogger.Println("C# project detected Ext sln")

				workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
				compile.CSharp(workspace)
				break
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

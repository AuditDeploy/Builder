package utils

import (
	"Builder/compile"
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//ProjectType will derive the poject type(go, node, java repo) and execute its compiler
func ProjectType() {

	//parentDir = the name of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	//languages we are currently compiling
	files := []string{"main.go", "package.json", "pom.xml"}

	for _, file := range files {

		filePath := parentDir + "/" + ".hidden" + "/" + file

		//checking if the filepath exists
		fileExists, err := exists(filePath)

		if err != nil {
			logger.ErrorLogger.Println("No Go, Npm, or Java File Exists")
			log.Fatal(err)
		}

		fmt.Println(fileExists)

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
			default:
				deriveProjectByExtension()
			}
		} else {
			fmt.Println("SHOULD NOT PRINT OUT")
			// deriveProjectByExtension()
		}
	}
}

//derive projects by Extensions
func deriveProjectByExtension() {
	//parentDir = the name of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	extensions := []string{".csproj"}

	for _, ext := range extensions {
		err := exec.Command("find", parentDir+"/"+".hidden", "-name", fmt.Sprintf("*%s", ext)).Run()

		if err != nil {
			log.Fatal(err)
		} else {
			switch ext {
			case ".csproj":
				CopyDir()
				logger.InfoLogger.Println("C# project detected")

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

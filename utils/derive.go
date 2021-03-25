package utils

import (
	"Builder/compile"
	"Builder/logger"
	"log"
	"os"
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
			}
		}
	}
}

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
package utils

import (
	"Builder/compile"
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//ProjectType will derive the project type and execute its compiler
func ProjectType() {

	//parentDir = the path of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")
	//check for user defined type from builder.yaml
	configType := strings.ToLower(os.Getenv("BUILDER_PROJECT_TYPE"))

	recurseExists()

	var files []string
	//projectType exists in builder.yaml
	if (configType != "") {
		//check value of config type, return string array of languages build file/files
		files = ConfigDerive()
	} else {
		//default
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
			if (file == "main.go" || configType == "go") {
					//executes go compiler
					CopyDir()
					logger.InfoLogger.Println("Go project detected")
					compile.Go(filePath)
				} else if (file == "package.json" || configType == "node" || configType == "npm") {
					//executes node compiler
					logger.InfoLogger.Println("Npm project detected")
					compile.Npm()
				} else if (file == "pom.xml" || configType == "java") {
					//executes java compiler
					CopyDir()
					logger.InfoLogger.Println("Java project detected")

					workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
					compile.Java(workspace)
				} else if (file == "gemfile.lock" || configType == "ruby") {
					//executes ruby compiler
					logger.InfoLogger.Println("Ruby project detected")
					compile.Ruby()
				} else if (file == "pipfile.lock" || configType == "python") {
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

func recurseExists() ([]string, error) {
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	fileList := make([]string, 0)
	e := filepath.Walk(hiddenDir, func(path string, f os.FileInfo, err error) error {
		if (f.Name() == "main.go") {
			fileList = append(fileList, path)
		}
		return err
	})
	
	if e != nil {
		panic(e)
	}

	for _, file := range fileList {
		fmt.Println(file)
	}

	return fileList, nil
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

package derive

import (
	"Builder/compile"
	"Builder/logger"
	"Builder/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		files = utils.ConfigDerive()
	} else {
		//default
		files = []string{"main.go", "package.json", "pom.xml", "gemfile.lock", "gemfile", "requirements.txt"}
	}

	var filePath string
	//look for those files inside hidden dir
	for _, file := range files {

		//recursively check for file in hidden dir, return path if found
		filePath = findPath(file)
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
				utils.CopyDir()
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
				finalPath := createFinalPath(filePath, file)

				utils.CopyDir()
				logger.InfoLogger.Println("Java project detected")

				compile.Java(finalPath)
				return
			} else if file == "gemfile.lock" || file == "gemfile" || configType == "ruby" {
				//executes ruby compiler
				logger.InfoLogger.Println("Ruby project detected")
				compile.Ruby()
				return
			} else if file == "requirements.txt" || configType == "python" {
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
	var dirPathExtToFound string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		path, _ := os.Getwd()
		dirPathExtToFound = path
	} else {
		dirPathExtToFound = os.Getenv("BUILDER_HIDDEN_DIR")
	}
	extensions := []string{".csproj", ".sln"}

	for _, ext := range extensions {
		extFound, fileName := extExistsFunction(dirPathExtToFound, ext)

		if extFound {
			switch ext {
			//checks if ext exists, if it's .csprocj it will pass down the filePath to c# compiler
			case ".csproj":
				filePath := strings.Replace(findPath(fileName), ".hidden", "workspace", 1)

				if os.Getenv("BUILDER_COMMAND") != "true" {
					utils.CopyDir()
				}
				logger.InfoLogger.Println("C# project detected, Ext .csproj")
				compile.CSharp(filePath)

			//if it's .sln, it will find all the project path in the solution(repo)
			case ".sln":
				filePath := strings.Replace(findPath(fileName), ".hidden", "workspace", 1)

				if os.Getenv("BUILDER_COMMAND") != "true" {
					utils.CopyDir()
				}

				listOfProjects, err := exec.Command("dotnet", "sln", filePath, "list").Output()

				if err != nil {
					log.Fatal(err)
				}

				stringifyListOfProjects := string(listOfProjects)
				listOfProjectsArray := strings.Split(stringifyListOfProjects, "\n")[2:]

				//if there's more than 5 projects in solution(repo), user will be asked to use builder config instead
				if len(listOfProjectsArray) > 5 {
					logger.InfoLogger.Println("C# project detected, Ext .sln. More than 5 projects in solution not supported")
					log.Fatal("There are more than 5 projects in this solution, please use Builder Config and specify the path of the file you wish to compile in the builder.yml")
				} else {
					var pathToCompileFrom string

					if os.Getenv("BUILDER_COMMAND") == "true" {
						// path, _ := os.Getwd()
						buildFile := os.Getenv("BUILDER_BUILD_FILE")
						pathToCompileFrom = buildFile
					} else {
						// < 5 projects in solution(repo), user will be prompt to choose a project path.
						pathToCompileFrom = selectPathToCompileFrom(listOfProjectsArray)
						workspace := os.Getenv("BUILDER_WORKSPACE_DIR")
						pathToCompileFrom = workspace + "/" + pathToCompileFrom
						utils.CopyDir()
					}

					logger.InfoLogger.Println("C# project detected, Ext .sln")
					compile.CSharp(pathToCompileFrom)
				}
			}
		}
	}

}

//takes in file, searches hiddenDir to find a match and returns path to file
func findPath(file string) string {

	var dirPath string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		currentDir, _ := os.Getwd()
		dirPath = currentDir
	} else {
		dirPath = os.Getenv("BUILDER_HIDDEN_DIR")
	}

	// if f.Name is == to file passed in "coolProject.go", filePath becomes the path that file exists in
	var filePath string
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
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

func extExistsFunction(dirPath string, ext string) (bool, string) {
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

	var fileName string

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ext {
				fileName = file.Name()
				found = true
			}
		}
	}

	return found, fileName
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

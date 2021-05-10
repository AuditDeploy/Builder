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

//derive the project type and execute its compiler
func ProjectType() {
	configType := strings.ToLower(os.Getenv("BUILDER_PROJECT_TYPE"))

	var filesToCompileFrom []string
	if configType != "" {
		filesToCompileFrom = utils.ConfigDerive()
	} else {
		filesToCompileFrom = []string{"main.go", "package.json", "pom.xml", "gemfile.lock", "gemfile", "requirements.txt"}
	}

	var filePath string
	for _, file := range filesToCompileFrom {

		filePath = findFilePathInHiddenDir(file)
		fileExists, err := fileExistsInDir(filePath)

		if err != nil {
			logger.ErrorLogger.Println("No Go, Npm, Ruby, Python or Java File Exists")
			log.Fatal(err)
		}
		//if file exists and filePath isn't empty, run conditional to find correct compiler
		if fileExists && filePath != "" && filePath != "./" {
			if file == "main.go" || configType == "go" {

				finalPath := createFinalPath(filePath, file)
				utils.CopyDir()
				logger.InfoLogger.Println("Go project detected")
				compile.Go(finalPath)
				return
			} else if file == "package.json" || configType == "node" || configType == "npm" {

				logger.InfoLogger.Println("Npm project detected")
				compile.Npm()
				return
			} else if file == "pom.xml" || configType == "java" {

				finalPath := createFinalPath(filePath, file)
				utils.CopyDir()
				logger.InfoLogger.Println("Java project detected")
				compile.Java(finalPath)
				return
			} else if file == "gemfile.lock" || file == "gemfile" || configType == "ruby" {

				logger.InfoLogger.Println("Ruby project detected")
				compile.Ruby()
				return
			} else if file == "requirements.txt" || configType == "python" {

				logger.InfoLogger.Println("Python project detected")
				compile.Python()
				return
			}
		}
	}

	//C# compiler
	deriveProjectByExtension()
}

//derive projects by Extensions
func deriveProjectByExtension() {

	var dirPathToFindExt string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		currentUserpPath, _ := os.Getwd()
		dirPathToFindExt = currentUserpPath
	} else {
		dirPathToFindExt = os.Getenv("BUILDER_HIDDEN_DIR")
	}

	extensions := []string{".csproj", ".sln"}

	for _, ext := range extensions {
		extFound, filePath := extExistsFunction(dirPathToFindExt, ext)

		if extFound {
			switch ext {
			case ".csproj":

				if os.Getenv("BUILDER_COMMAND") != "true" {
					utils.CopyDir()
				}

				logger.InfoLogger.Println("C# project detected, Ext .csproj")
				compile.CSharp(filePath)

			case ".sln":

				if os.Getenv("BUILDER_COMMAND") != "true" {
					utils.CopyDir()
				}

				listOfProjectsArray := ListAllProjectsInSolution(filePath)

				//if there's more than 5 projects in solution(repo), user will be asked to use builder config instead
				if len(listOfProjectsArray) > 5 {

					logger.InfoLogger.Println("C# project detected, Ext .sln. More than 5 projects in solution not supported")
					log.Fatal("There are more than 5 projects in this solution, please use Builder Config and specify the path of the file you wish to compile in the builder.yml")

				} else {
					var pathToCompileFrom string

					if os.Getenv("BUILDER_COMMAND") == "true" {

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

func ListAllProjectsInSolution(filePath string) []string {
	listOfProjects, err := exec.Command("dotnet", "sln", filePath, "list").Output()

	if err != nil {
		log.Fatal(err)
	}

	stringifyListOfProjects := string(listOfProjects)
	listOfProjectsArray := strings.Split(stringifyListOfProjects, "\n")[2:]

	return listOfProjectsArray
}

func findFilePathInHiddenDir(file string) string {

	var dirPath string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		currentDir, _ := os.Getwd()
		dirPath = currentDir
	} else {
		dirPath = os.Getenv("BUILDER_HIDDEN_DIR")
	}

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

func fileExistsInDir(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func extExistsFunction(dirPath string, ext string) (bool, string) {
	extFound := false

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
				extFound = true
			}
		}
	}

	var filePath string
	if fileName != "" {
		filePath = strings.Replace(findFilePathInHiddenDir(fileName), ".hidden", "workspace", 1)
	}

	return extFound, filePath
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

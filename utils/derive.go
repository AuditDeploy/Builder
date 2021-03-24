package utils

import (
	// "fmt"
	"log"
	"os"
	// "os/exec"
	// "strings"

	compile "Builder/compile"
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
			log.Fatal(err)
		}

		//if file exists run a swith statement
		if fileExists {
			switch file {
			case "main.go":
				//executes go compiler
				CopyDir()
				compile.Go(filePath)
			case "package.json":
				//executes node compiler
				compile.Npm()
			case "pom.xml":
				//executes java compiler
				CopyDir() 
				// // filePath2, _ := exec.Command("find", parentDir+"/workspace", "-name", "*.xml").CombinedOutput()

				// //returning multiple possible paths which are separated by a newline "\n"
				// stringPath := string(filePath2)

				// //split paths are returns an array of paths
				// paths := strings.Split(stringPath, "\n")
				// // fmt.Println(paths[0])
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
package derive

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	compile "github.com/ilarocca/Builder/compile"
)

//ProjectType will derive the poject type(go, node, java repo) and execute its compiler
func ProjectType() {

	//parentDir = the name of the project
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	//languages we are currently compiling
	files := []string{"main.go", "package.json", "pom.xml"}

	for _, file := range files {

		filePath := parentDir + "/" + "workspace" + "/" + file

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
				compile.Go(filePath)
			case "package.json":
				//executes node compiler
				fmt.Println("FILEEEE package.json")
			case "pom.xml":
				//executes java compiler
				filePath2, _ := exec.Command("find", "./helloworld/workspace", "-name", "*.java").CombinedOutput()

				//returning multiple possible paths which are separated by a newline "\n"
				stringPath := string(filePath2)

				//split paths are returns an array of paths
				paths := strings.Split(stringPath, "\n")

				compile.Java(paths[0])
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

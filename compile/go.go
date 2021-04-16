// takes in code as arg from go
//run go build on code given

package compile

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Go creates exe from file passed in as arg
func Go(filePath string) {

	//define dir path for command to run in
	var fullPath string
	configPath := os.Getenv("BUILDER_DIR_PATH")
	//if user defined path in builder.yaml, full path is included already, else add curren dir + local path 
	if (configPath != "") {
		// ex: C:/Users/Name/Projects/helloworld_19293/workspace/dir
		fullPath = filePath
	} else {
		path, _ := os.Getwd()
		//combine local path to newly created tempWorkspace, 
		//gets rid of "." in path name
		// ex: C:/Users/Name/Projects + /helloworld_19293/workspace/dir
		fullPath = path + filePath[strings.Index(filePath, ".")+1:]
	}
	
	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	//find 'go file' to be built
	buildFile := strings.ToLower(os.Getenv("BUILDER_BUILD_FILE"))
	//if no file defined by user, use default main.go
	if (buildFile == "") {
		buildFile = "main.go"
	}

	var cmd *exec.Cmd
	if (buildTool == "go") {
		fmt.Println(buildTool)
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath       // or whatever directory it's in
	} else {
		//default
		cmd = exec.Command("go", "build", buildFile)
		cmd.Dir = fullPath       // or whatever directory it's in
	}

	fmt.Println("go full path: "+cmd.Dir)

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Go project failed to compile.")
		log.Fatal(err)
	}

	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", fullPath+"/main.exe", artifactPath).Run()
	}

	logger.InfoLogger.Println("Go project compiled successfully.")
}


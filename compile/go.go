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
func Go(filepath string) {

	//copies contents of .hidden to workspace
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	var cmd *exec.Cmd
	if (buildTool == "go") {
		fmt.Println(buildTool)
		cmd = exec.Command("go", "build", "-o", workspaceDir, filepath)
	} else {
		//default
		cmd = exec.Command("go", "build", "-o", workspaceDir, filepath)
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Go project failed to compile.")
		log.Fatal(err)
	}

	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", workspaceDir+"/main.exe", artifactPath).Run()
	}

	logger.InfoLogger.Println("Go project compiled successfully.")
}


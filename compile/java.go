package compile

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Java does ...
func Java(filePath string) {

	//install dependencies/build, if yaml build type exists install accordingly
	buildType := strings.ToLower(os.Getenv("BUILDER_BUILD_TYPE"))
	var cmd *exec.Cmd
	if (buildType == "maven" || buildType == "mvn") {
		fmt.Println(buildType)
		cmd = exec.Command("mvn", "clean", "install", "-f", filePath)
	} else if (buildType == "grendle") {
		// grendle, etc.
	} else {
		//default
		cmd = exec.Command("mvn", "clean", "install", "-f", filePath)
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Java project failed to compile.")
		log.Fatal(err)
	}

	//if artifact path exists, copy contents
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", workspaceDir+"/target", artifactPath).Run()
	}
	logger.InfoLogger.Println("Java project compiled successfully.")

}
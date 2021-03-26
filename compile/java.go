package compile

import (
	"Builder/logger"
	"log"
	"os"
	"os/exec"
)

//Java does ...
func Java(filePath string) {

	//copies contents of .hidden to workspace
	cmd := exec.Command("mvn", "clean", "install", "-f", filePath)
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	

	if err != nil {
		logger.ErrorLogger.Println("Java project failed to compile.")
		log.Fatal(err)
	}

	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", workspaceDir+"/target", artifactPath).Run()
	}
	logger.InfoLogger.Println("Java project compiled successfully.")

}
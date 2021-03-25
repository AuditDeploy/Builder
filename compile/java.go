package compile

import (
	"Builder/logger"
	"log"
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

	logger.InfoLogger.Println("Java project compiled successfully.")

}
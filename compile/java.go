package compile

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//Java does ...
func Java(filePath string) {

	//copies contents of .hidden to workspace
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	fmt.Println(filePath)
	fmt.Println(workspaceDir)


	cmd := exec.Command("javac", filePath)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	//make contents read-only
	exec.Command("chmod", "-R", "0444", hiddenDir).Run()
}
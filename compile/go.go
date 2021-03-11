// takes in code as arg from go
//run go build on code given

package compile

import (
	"log"
	"os"
	"os/exec"
)

//Go creates exe from file passed in as arg
func Go() {

	//copies contents of .hidden to workspace
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	exec.Command("cp", "-a", hiddenDir+"/.", workspaceDir).Run()

	cmd := exec.Command("go", "mod", "init")
	cmd.Run()

	//compile source code in workspace
	cmd2 := exec.Command("go", "build", "-o", workspaceDir, "main.go")
	err := cmd2.Run()

	if err != nil {
		log.Fatal(err)
	}

	//make contents read-only
	exec.Command("chmod", "-R", "0444", hiddenDir).Run()
}


// takes in code as arg from go
//run go build on code given

package compile

import (
	"fmt"
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

	//compile source code in workspace
	cmd := exec.Command("go", "build", "-o", workspaceDir, "main.go")

	//search for a 'main.go' filename and add that path to workspaceDir
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(stdout))
}


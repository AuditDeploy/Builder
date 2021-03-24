// takes in code as arg from go
//run go build on code given

package compile

import (
	"log"
	"os"
	"os/exec"
)

//Go creates exe from file passed in as arg
func Go(filepath string) {

	//copies contents of .hidden to workspace
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	cmd := exec.Command("go", "mod", "init")
	cmd.Run()

	//compile source code in workspace
	cmd2 := exec.Command("go", "build", "-o", workspaceDir, filepath)
	err := cmd2.Run()

	if err != nil {
		log.Fatal(err)
	}

}


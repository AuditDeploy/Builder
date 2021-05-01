// takes in code as arg from go
//run go build on code given

package dependencies

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//Go creates exe from file passed in as arg
func Go() {

	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	cmd := exec.Command("go", "mod", "init", workspaceDir)

	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(stdout))
}

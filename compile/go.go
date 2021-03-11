// takes in code as arg from go
//run go build on code given

package compile

import (
	"os"
	"os/exec"
)

//Go creates exe from file passed in as arg
func Go(fileToCompile string) {

	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")

	cmd := exec.Command("go", "build", "-o", workspaceDir, fileToCompile)
	cmd.Run()

	//search for a 'main.go' filename and add that path to workspaceDir
	// stdout, err := cmd.Output()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Print(string(stdout))

	// Artifact2()
}

// //Artifact2 does ...
// func Artifact2() {
// 	hidden := os.Getenv("BUILDER_HIDDEN_DIR")
// 	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
// 	exec.Command("cp", workspaceDir+"/main", hidden).Run()
// }

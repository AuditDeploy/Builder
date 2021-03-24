package compile

import (
	"log"
	"os/exec"
)

//Java does ...
func Java(filePath string) {

	//copies contents of .hidden to workspace
	cmd := exec.Command("mvn", "clean", "install", "-f", filePath)
	err := cmd.Run()
	

	if err != nil {
		log.Fatal(err)
	}
}
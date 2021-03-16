package compile

import (
	"log"
	"os"
	"os/exec"
)

//Java does ...
func Java(filePath string) {

	//copies contents of .hidden to workspace
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	cmd := exec.Command("javac", filePath)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	//make hiddenDir hidden
	exec.Command("attrib", hiddenDir, "-h").Run()
	//make contents read-only
	exec.Command("chmod", "-R", "0444", hiddenDir).Run()
}
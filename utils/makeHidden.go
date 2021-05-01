package utils

import (
	"os"
	"os/exec"
)

func MakeHidden() {
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")

	//make hiddenDir hidden
	exec.Command("attrib", hiddenDir, "-h").Run()
	//make contents read-only
	exec.Command("chmod", "-R", "0444", hiddenDir).Run()
}

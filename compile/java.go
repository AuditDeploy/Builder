package compile

import (
	"fmt"
	"log"
	"os/exec"
)

//Java does ...
func Java(filePath string) {

	fmt.Println("compiler start")

	cmd := exec.Command("javac", filePath)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("compiler End")
}

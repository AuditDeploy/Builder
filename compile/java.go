package compile

import (
	"fmt"
	"log"
	"os/exec"
)

//Java does ...
func Java() {

	fmt.Println("compiler start")
	cmd := exec.Command("javac", "javaHelloWorld/HelloWorld.java")
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("compiler End")
}

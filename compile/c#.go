package compile

import (
	"log"
	"os/exec"
)

func CSharp(filepath string) {
	err := exec.Command("dotnet", "build", filepath).Run()

	if err != nil {
		log.Fatal(err)
	}
}

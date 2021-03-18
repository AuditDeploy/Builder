package cmd

import (
	"log"
	"os"
	"os/exec"
)

func InitCmd() {
	init := exec.Command("go", "run", "main.go")
	init.Run()

	if err != nil {
		log.Fatal(err)
	}
}
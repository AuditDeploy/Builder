package main

import (
	"Builder/cmd"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) > 1 {
		builderCommand := os.Args[1]

		if builderCommand == "init" {
			cmd.Init()
		} else if builderCommand == "config" {
			cmd.Config()
		} else {
			fmt.Println("expected command: 'init' or 'config'")
		}
	} else {
		builder()
	}
}

func builder() {
	path, _ := os.Getwd()
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		exec.Command("git", "pull").Run()

		//project type
		//buildFile
		//buildPath
		exec.Command("go", "build", path+"/"+"main.go").Run()
	} else {
		log.Fatal("bulder.yaml file not found or cd into workspace")
	}

}

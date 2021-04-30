package main

import (
	"Builder/cmd"
	"Builder/derive"
	"Builder/logger"
	"Builder/yaml"
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
		cmd.Builder()
	}
	fmt.Println("Build Complete 🔨")
}
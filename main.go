package main

import (
	"builder/cmd"
	"fmt"
	"os"
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
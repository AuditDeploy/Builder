package main

import (
	"builder/cmd"
	"builder/utils"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		utils.Help()
		builderCommand := os.Args[1]
		if builderCommand == "init" {
			cmd.Init()
		} else if builderCommand == "config" {
			cmd.Config()
		} else {
			cmd.Builder()
		}
	} else {
		cmd.Builder()
	}
	fmt.Println("Build Complete 🔨")
}

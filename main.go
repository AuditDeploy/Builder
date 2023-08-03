package main

import (
	"Builder/cmd"
	"Builder/gui"
	"Builder/utils"
	"Builder/utils/log"
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
		} else if builderCommand == "gui" {
			gui.Gui()
		} else {
			cmd.Builder()
		}
	} else {
		cmd.Builder()
	}

	log.Info("Build Complete 🔨")
	fmt.Println("Build Complete 🔨")
}

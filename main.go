package main

import (
	"Builder/cmd"
	"Builder/gui"
	"Builder/utils"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var BuilderLog = zap.S()

func main() {

	if len(os.Args) > 1 {
		utils.Help()
		builderCommand := os.Args[1]
		if builderCommand == "init" {
			cmd.Init()
			fmt.Println("Build Complete 🔨")
		} else if builderCommand == "config" {
			cmd.Config()
			fmt.Println("Build Complete 🔨")
		} else if builderCommand == "gui" {
			gui.Gui()
		} else {
			cmd.Builder()
			fmt.Println("Build Complete 🔨")
		}
	} else {
		cmd.Builder()
		fmt.Println("Build Complete 🔨")
	}
}

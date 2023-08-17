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

	BuilderLog.Info("Build Complete 🔨")
	fmt.Println("Build Complete 🔨")
}

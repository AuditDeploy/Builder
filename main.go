package main

import (
	"Builder/cmd"
	"Builder/logger"
	"fmt"
	"os"
)

func main() {
		logger.InfoLogger.Println("Builder is starting...")
	
		builderCommand := os.Args[1]

		if builderCommand == "init" {
			cmd.Init()
		} else if builderCommand == "config" {
			cmd.Config()
		} else {fmt.Println("expected command: 'init' or 'config'")}

}

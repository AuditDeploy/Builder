package main

import (
	"Builder/cmd"
	"fmt"
	"os"
)

func main() {	
		builderCommand := os.Args[1]

		if builderCommand == "init" {
			cmd.Init()
		} else if builderCommand == "config" {
			cmd.Config()
		} else {fmt.Println("expected command: 'init' or 'config'")}
}

package cmd

import (
	"Builder/utils"
	"fmt"
)

func Config() {
	fmt.Println("this is the 'config' command")
	utils.YamlParser()
}
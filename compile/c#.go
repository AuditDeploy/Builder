package compile

import (
	"Builder/logger"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CSharp(filepath string) {

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	var cmd *exec.Cmd
	if buildTool == "dotnet" {
		fmt.Println(buildTool)
		cmd = exec.Command("dotnet", "build", filepath)

	} else {
		//default
		cmd = exec.Command("dotnet", "build", filepath)
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("C# project failed to compile.")
		log.Fatal(err)
	}

	logger.InfoLogger.Println("C# project compiled successfully.")
}

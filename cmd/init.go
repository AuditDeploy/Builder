package cmd

import (
	
	directory "github.com/ilarocca/Builder/directory"
	utils "github.com/ilarocca/Builder/utils"
)

func Init() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	// make dirs
	directory.MakeParentDir()

	// clone repo into hidden
	utils.CloneRepo()

	// compile logic to derive project type
	utils.ProjectType()
	// copy hidden into work dir, install dependencies, compile source code from repo

}
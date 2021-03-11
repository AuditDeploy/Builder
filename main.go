package main

import (
	"os"

	compile "github.com/ilarocca/Builder/compile"
	directory "github.com/ilarocca/Builder/directory"
	utils "github.com/ilarocca/Builder/utils"
)

func main() {
	//check argument syntax, exit if incorrect
	utils.CheckArgs()

	args := os.Args[1:]
	repoURL := os.Args[2]

	//check args/flags
	// clone repo

	// make dirs
	directory.MakeParentDir(repoURL)

	// clone repo into hidden
	utils.CloneRepo(args)

	// compile logic to derive project type 

	// copy hidden into work dir, install dependencies, compile source code from repo
	compile.Npm()
}


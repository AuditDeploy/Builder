package main

import (
	"os"

	compile "github.com/ilarocca/Builder/compile"
	directory "github.com/ilarocca/Builder/directory"
	utils "github.com/ilarocca/Builder/utils"
)

func main() {
	args := os.Args[1:]
	repoURL := os.Args[2]

	// make dirs
	directory.MakeParentDir(repoURL)

	// clone repo into hidden
	utils.CloneRepo(args)

	//install dependecies ('npm install', etc)
	// dependency.Go()

	//QUESTIONS install dependencies in compile package?
	//make hidden after compile? 

	// compile logic to derive project type

	// compile source code from repo
	compile.Npm()

	// pass source code into hidden dir AND workspace

}
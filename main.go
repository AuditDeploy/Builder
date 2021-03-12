package main

import (
	"os"

	compile "github.com/ilarocca/Builder/compile"
	parent "github.com/ilarocca/Builder/makeDirFunctions/parent"
	cloneRepo "github.com/ilarocca/Builder/utils/repo"
)

func main() {
	args := os.Args[1:]
	// repoURL := os.Args[2]

	// make parent
	// make sub dirs
	parent.MakeParentDir()

	// clone repo into workspace
	cloneRepo.GetURL(args)

	//install dependecies ('npm install', etc)
	// dependency.Go()

	// compile logic to derive project type
	// compile source code from repo
	compile.Go()

	// pass source code into hidden dir AND workspace

}

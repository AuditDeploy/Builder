package main

import (
	"os"

	parent "github.com/ilarocca/Builder/makeDirFunctions/parent"
	derive "github.com/ilarocca/Builder/utils/derive"
	cloneRepo "github.com/ilarocca/Builder/utils/repo"
)

func main() {
	args := os.Args[1:]
	repoURL := os.Args[2]

	// make parent
	// make sub dirs
	parent.MakeParentDir(repoURL)

	// clone repo into workspace
	cloneRepo.GetURL(args)

	//derive project type and compiles repo
	derive.ProjectType()

	//install dependecies ('npm install', etc)
	// dependency.Go()

	// compile logic to derive project type
	// compile source code from repo
	// compile.Go()

	// pass source code into hidden dir AND workspace

}

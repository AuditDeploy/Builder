package utils

import (
	"fmt"
	"os"
)

//Help is application info
func Help() {
	cArgs := os.Args[1:]
	var helpExists bool

	for _, v := range cArgs {
		if v == "--help" || v == "-h" {
			helpExists = true
		}
	}

	//check for help flag or builder cmd to print info
	if (os.Getenv("BUILDER_COMMAND") == "true") || helpExists {
		fmt.Println(`
		   🔨 BUILDER 🔨
													
	       #%&&&%  ,&&            
	    ##. #&&&&&&&&& &&&&&      
		.&&&#        &&&&/    
		.&&&%         &&&&    
		.&&&#        &&&&,    
		.&&&% &&&&&&&&&&      
		.&&&# ......,#&&&&%   
		.&&&#           &&&&  
		.&&&#           #&&&. 
		.&&&#          %&&&#  
		.&&&% &&&&&&&&&&&&.   
		.&&&% &&&&&&&#,       										

			Commands

* builder init: auto-build a project that doesn't have a builder yaml (repo needed)
	- ex: builder init <repo> <flags> 
* builder config: build project w/ a builder.yaml (repo needed)
	- ex: builder config <repo> <flags>
* builder: build project w/ builder.yaml while in the projects directory (no repo needed) 
	- ex: builder <flags> 

			Flags

* '--help' or '-h': provide info for Builder
* '--output' or '-o': user defined output path for artifact
* '--name' or '-n': user defined project name
* '--yes' or '-y': bypass prompts
* '--branch' or '-b': specify repo branch


		builder.yaml params
* projectName: provide name for project
  - ("helloworld", etc)
* projectPath: provide path for project to be built
  - ("/Users/Name/Projects", etc)
* projectType: provide language/framework being used
  - (Node, Java, Go, Ruby, Python, C#, Ruby)
* buildTool: provide tool used to install dependencies/build project
  - (maven, npm, bundler, pipenv, etc)
* buildFile: provide file name needed to install dep/build project
  - Can be any user specified file. (myCoolProject.go, package.json etc)
* buildCmd: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)
* outputPath: provide path for artifact to be sent
  - ("/Users/Name/Artifacts", etc)
* globalLogs: specify path to global logs
  - ("var/logs/global-logs/logs.txt")
* dockerCmd: specify docker command, if building a container
  - ("docker build -t my-project:1.3 .")
* repoBranch: specify repo branch name
  - (“feature/“new-branch”)
* bypassPrompts: bypass prompts
  - (true)
			`)
		os.Exit(0)
	}
}

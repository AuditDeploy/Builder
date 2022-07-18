package utils

import (
	"fmt"
	"os"
)

//CheckArgs is...
func Help() {
	cArgs := os.Args[1:]

	//check for help flag
	for _, v := range cArgs {
		if v == "--help" || v == "-h" {
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

* projectpath: provide path for project to be built
  - ("/Users/Name/Projects", etc)
* projecttype: provide language/framework being used
  - (Node, Java, Go, Ruby, Python, C#, Ruby)
* buildtool: provide tool used to install dependencies/build project
  - (maven, npm, bundler, pipenv, etc)
* buildfile: provide file name needed to install dep/build project
  - Can be any user specified file. (myCoolProject.go, package.json etc)
* buildcmd: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)
* outputpath: provide path for artifact to be sent
  - ("/Users/Name/Artifacts", etc)
* globallogs: specify path to global logs
  - ("var/logs/global-logs/logs.txt")
* dockercmd: specify docker command, if building a container
  - ("docker build -t my-project:1.3 .")
			`)
			os.Exit(0)
		}
	}
}

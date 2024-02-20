package utils

import (
	"fmt"
	"os"
)

// Help is application info
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
* builder push [optional flags]: pushes build metadata and logs JSON to provided url (url needed in line or in builder.yaml)
    - optional flag: '--save' to automate "push" process for future builds
	- ex: builder push <optional_flags> <url>
* builder: build project w/ builder.yaml while in the projects directory (no repo needed) 
	- ex: builder <flags> 
* builder gui: display the Builder GUI (requires Chrome for use)

			Flags

* '--help' or '-h': provide info for Builder
* '--output' or '-o': user defined output path for artifact
* '--name' or '-n': user defined project name
* '--branch' or '-b': specify repo branch
* '--debug' or '-d': show Builder log output
* '--verbose' or '-v': show log output for project being built
* '--docker' or '-D': build Docker image


		builder.yaml params
		
* projectname: provide name for project
  - ("helloworld", etc)
* projectpath: provide path for project to be built
  - ("/Users/Name/Projects", etc)
* projecttype: provide language/framework being used
  - ("Node", "Java", "Go", "Rust", "Python", "C#", "Ruby")
* buildtool: provide tool used to install dependencies/build project
  - ("maven", "npm", "bundler", "pipenv", etc)
* buildfile: provide file name needed to install dep/build project
  - Can be any user specified file. ("myCoolProject.go", "package.json" etc)
* prebuildcmd: for C/C++ projects only. Provide command to run before configcmd and buildcmd 
  - ("autoreconf -vfi", "./autogen.sh", etc)
* configcmd: for C/C++ projects only. provide full command to configure C/C++ project before running buildcmd
  - ("./configure")
* buildcmd: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)
* artifactlist: provide comma seperated list of artifact names as string
  - ("artifact", "artifact.exe", "artifact.rpm,artifact2.rpm,artifact3.rpm", etc)
* outputpath: provide path for artifact to be sent
  - ("/Users/Name/Artifacts", etc)
* repobranch: specify repo branch name
  - (“feature/“new-branch”)
* docker: generate docker image
  - dockerfile: name of dockerfile
  - registry: registry to push docker image to
  - version: tag to give docker image
* push:
  - url: url to push build metadata and logs JSON to
  - auto: (true/false) whether to automate pushing process for future builds
* appicon: specify url to app icon image
  - ("http://domain.co/path/to/app_icon.png")
		`)
		os.Exit(0)
	}
}

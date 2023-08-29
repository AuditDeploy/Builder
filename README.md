![alt text](https://uploads-ssl.webflow.com/63eabcf563ffc28443d6b895/64555593669ad6181f9aabba_AD_Builder-logo_orange%20(1).png)

# Builder

> A build tool with transparent logs and metadata tied to each build.

The only build tool that allows you to pass in a repo and build your project, while getting extensive metadata tied to your shippable product. There's no need to frantically track down build info if something goes wrong. Know who built it, in what directory, on what machine, at what time. Never not know again.

## Get Started

Once you pull the project down, you need to create an executable and place it in your the correct path. 

To create a binary run `go build -o builder`
Then move builder to a path of your choosing, e.g. `mv builder /usr/local/bin`

## About

Run `builder init [repo url]` on a project (must be a compatible language). This creates a new project with a hidden, logs, workspace, and artifact dir.

For compatible compiled languages (C#, Golang, Java) a builder.yaml file gets created (see below for more info) in the workspace dir. For interpreted languages (Javascript, Python, Ruby) it's created in the temp dir.

For compiled, the artifact dir contains an executable file, both json and yaml metadata, and a zip of those contents. That zip can be sent off wherever by specifing the "--output" path flag.

You can then cd into the path of your builder.yaml and run the `builder` command. This will pull any changes from the repo and create a new artifact (meaning a shippable product, either zip or executable) with new metadata.

IF you know your project has a specific buildFile name or you would like to use your own buildCommand, etc, you can attach a builder.yaml to your project and run `builder config [repo url]` instead of init. This config command allows for an extra layer of custom parameters and will use your builder.yaml instead of creating one.

Important to note: at this time, if you create your own builder.yaml and initialize a project with config, you should include all of the builder parameters (even if they are empty) in order for the `builder` command to work properly. This will be address in the next version.

Important note for Windows users: at this time, please use Git Bash when running any Builder command

## Builder CLI Exec Commands & Flags

### Commands:

Builder is great at guessing what to do with most repos it's given, for the other case, you need to initialize your project by placing a user defined builder.yaml in your repo and running the `config` command.

- `builder init`: auto default build a project (creates packaged artifact) with metadata and logs, creates default builder.yaml
  - only necessary argument is a github repo url
- `builder config`: user defined (user created builder.yaml) project build that creates artifact with metadata and logs
  - only necessary argument is a github repo url
- `builder`: user cds into a project path with a builder.yaml, it then pulls changes, creates new artifact and new metadata
  - no arguments accepted at this time
  - if you would like the new artifact sent to a specified dir, make sure your output path is specified in the builder.yaml
- `builder gui`: display the Builder GUI.  Requires Chrome for use

### Flags:

- '--help' or '-h': provide info for Builder
- '--output' or '-o': user defined output path for artifact
- '--name' or '-n': user defined project name
- '--branch' or '-b': specify repo branch
- '--debug' or '-d': show Builder log output
- '--verbose' or '-v': show log output for project being built
- '--docker' or '-D': build Docker image

## Builder Compatibility

### Languages/Frameworks with default build/install commands:

You must have the language or package manager previously installed in order to build specified project.

- Golang
  - Uses `go build main.go` as default command.
  - Uses `main.go` as entry point to project by default.
  - If your main package has a different name than main.go you need to create a builder.yaml within your repo, specify the buildfile, and run the `config` command.
- Node
  - Uses `npm install` as default command
  - Must have package.json in order to install dependencies by default.
- Java
  - Uses `mvn clean install` as default command.
  - Must have pom.xml as default buildfile.
- C#
  - Uses `dotnet build [file path]` as default command.
- Python
  - Uses `pip3 install -r requirements.txt -t [path/requirements]` as default command.
  - As of now, a requirements.txt is necessary to build default python projects.
- Ruby
  - Uses `bundle install --path vendor/bundle` as default command.
- C/C++
  - Looks for `Makefile` and runs `make` as default command.
  - To run autotools or a `./configure` command please specify these in the builder.yaml

To use other buildtools, buildcommands, or custome buildfiles you must create builder.yaml and run `config`.

## Builder.yaml Parameters

If you are specifying a buildfile, buildtool, or buildcmd within the builder.yaml, you MUST include the projectType.

At this point in time, please include ALL builder.yaml parameters (all keys must be lowercase), even if they are empty. (This will be addressed in the next update)

- `projectname`: provide name for project
  - ("helloworld", etc)
- `projectpath`: provide path for project to be built
  - ("/Users/Name/Projects", etc)
- `projecttype`: provide language/framework being used
  - ("Node", "Java", "Go", "Ruby", "Python", "C#", "Ruby", "C", "C++")
- `buildtool`: provide tool used to install dependencies/build project
  - ("maven", "npm", "bundler", "pipenv", etc)
  - for C/C++ project, please provide a build specific build tool from the following:
    - "make-rpm", "make-deb", "make-tar", "make-lib", "make-dll", or default "make" to build .exe files
- `buildfile`: provide file name needed to install dep/build project
  - Can be any user specified file. ("myCoolProject.go", "package.json", etc)
- `prebuildcmd`: for C/C++ projects only.  Provide command to run before configcmd and buildcmd
  - ("autoreconf -vfi", "./autogen.sh", etc)
- `configcmd`: for C/C++ projects only. provide full command to configure C/C++ project before running buildcmd
  - ("./configure")
- `buildcmd`: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)
- `artifactlist`: provide comma seperated list of artifact names as string
  - ("artifact", "artifact.exe", "artifact.rpm,artifact2.rpm,artifact3.rpm", etc)
- `outputpath`: provide path for artifact to be sent.  Please put the path in single quotes (')
  - ('/Users/Name/Artifacts', 'C:\Users\Name\Artifacts' etc)
- `dockercmd`: specify docker command, if building a container
  - ("docker build -t my-project:1.3 .")
- `repobranch`: specify repo branch name
  - (“feature/“new-branch”)

## Builder ENV Vars

### Native env vars:

- "BUILDER_PARENT_DIR": parent dir path
- "BUILDER_ARTIFACT_DIR": parent dir path
- "BUILDER_HIDDEN_DIR": hidden dir path
- "BUILDER_LOGS_DIR": logs dir path
- "BUILDER_COMMAND": bool if builder cmd is running

### Envs set by builder.config:

- "BUILDER_DIR_PATH": user defined parent dir path for specific build
- "BUILDER_PROJECT_TYPE": user defined project type ("go", "java", etc)
- "BUILDER_BUILD_TOOL": user defined build tool ("maven", "gradle", "npm", "yarn", etc)
- "BUILDER_BUILD_FILE": user defined build file ("myCoolProject.go")
- "BUILDER_BUILD_COMMAND": user defined build commmand ("yarn install")
- "BUILDER_ARTIFACT_LIST": user defined list of produced artifacts ("myProject.exe", "artifact.rpm,artifact2.rpm", etc)
- "BUILDER_OUTPUT_PATH": user defined output path for artifact

## Builder Funcionalty Layout

### main.go:

- check for either 'init' or 'config' command

### init:

#### 1. CheckArgs:

- checks arguments passed into the cli exec call
- call GetRepoURL
- check for '--help or -h' flag
- if no repo, exit
- if repo exists, check to ls-remote to see if it's a real git repo
- check for '--output or -o' flag (artifact/output path)
  - set 'BUILDER_OUTPUT_PATH' (either "" or user defined path)

#### 2. MakeDirs:

- call GetName (checks for '-n' flag, assigns dir name to name var)
- create 'configPath'var based on 'BUILDER_PATH_DIR' env var (established in a builder.yaml)
- create 'path' var either locally or with configPath + name + timestamp
- call MakeParentDir:
  - check if path already exists
  - make entire parentDir path
  - set 'BUILDER_PARENT_DIR' env var
- call MakeHiddenDir:
  - create hiddenPath var (parentDir + './hidden')
  - check if path already exists
  - make entire hiddenDir path
  - set 'BUILDER_HIDDEN_DIR' env var
- call MakeLogDir:
  - create logPath var (parentDir + '/logs')
  - check if path already exists
  - make entire logDir path
  - set 'BUILDER_LOGS_DIR' env var
  - call CreateLogs:
    - grab parentDir to create logs.txt name
    - create logs.txt in logs dir
    - create InfoLogger, WarningLogger, ErrorLogger vars to be used throughout
    - **_Logs.txt created HERE_**
- call MakeWorkspaceDir:
  - create workPath var (parentDir + '/workspace')
  - check if path already exists
  - make entire wokspaceDir path
  - set 'BUILDER_WORKSPACE_DIR' env var

#### 3. CloneRepo:

- call GetRepoURL:
  - check for repo after 'init' or 'config'
  - return repo string
- check "BUILDER_HIDDEN_DIR', if "", clone into tempRepo dir (used for 'config')
- clone repo into hiddenDir

#### 4. ProjectType:

- check "BUILDER_PROJECT_TYPE", if exists call ConfigDerive:
  - check "BUILDER_BUILD_FILE", if exists, return user specified build file/files
  - checks env var, returns []string containing languages default build file/files
- set files []string to builder.yaml val or default
- cycle through hidden dir to find the project type
- if file path for one of the project types exists, compile project
- GO -->
  - copy contents of hidden into workspace dir
  - compile.Go:
    - check "BUILDER_BUILD_TOOL" if exists, run that build tool, else run default
    - run 'go build' (default) in workspace path
    - if "BUILDER_OUTPUT_PATH" exists, copy artifact to that path
- JAVA -->
  - copy contents of hidden into workspace dir
  - compile.Java:
    - check "BUILDER_BUILD_TOOL" if exists, run that build tool, else run default
    - run 'mvn clean install' (default) in workspace path
    - if "BUILDER_OUTPUT_PATH" exists, copy artifact to that path
- NPM -->
  - compile.Npm:
    - create temp directory inside workspace dir
    - copy hidden dir contents (repo) into temp dir
    - check "BUILDER_BUILD_TOOL" if exists, run that build tool, else run default
    - run 'npm install' (default) in temp dir path
    - create temp.zip dir
    - recursively add files from temp dir to temp.zip
    - if "BUILDER_OUTPUT_PATH" exists, copy artifact (zip file in this case) to that path

#### 5. Metadata:

- create a yaml & json inside parent dir with:
  - ProjectName
	- ProjectType
	- ArtifactName
	- ArtifactChecksum
	- ArtifactLocation
	- UserName
	- HomeDir
	- IP
	- StartTime
	- EndTime
	- GitURL
	- MasterGitHash
	- BranchName

#### 6. MakeHidden:

- give hidden dir the hidden attrib
- make hidden dir read-only

#### 7. GlobalLogs:

- create globalLogs dir and .txt if it doesn't exists
- copy current builds logs.txt contents into globalLogs.txt

### config:

#### 1. CheckArgs:

- checks arguments passed into the cli exec call
- call GetRepoURL
- check for '--help or -h' flag
- if no repo, exit
- if repo exists, check to ls-remote to see if it's a real git repo
- check for '--output or -o' flag (artifact/output path)
  - set 'BUILDER_OUTPUT_PATH' (either "" or user defined path)

#### 2. CloneRepo:

**_This instance is to clone the repo in tempRepo dir to get builder.yaml info_**

- call GetRepoURL:
  - check for repo after 'init' or 'config'
  - return repo string
- check "BUILDER_HIDDEN_DIR', if "", clone into tempRepo dir (used for 'config')
- clone repo into hiddenDir

#### 3. YamlParser:

- create a map[string] interface{} to dump yaml info into
- read builder.yaml in the tempRepo dir
- unpack the yaml file into the map int{}
- pass the map int{} into ConfigEvens:
  - check for specific keys in map and create env vars based on value
  - check "projectType" and create 'BUILDER_PROJECT_TYPE' env var
  - check "buildTool" and create 'BUILDER_BUILD_TOOL' env var
  - check "buildFile" and create 'BUILDER_BUILD_FILE' env var
  - check "path" (this is parent dir path) and create 'BUILDER_DIR_PATH' env var
- delete tempRepo dir

#### 4. Run same functionality as 'init'

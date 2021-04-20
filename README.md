# Builder OS

## Builder CLI Exec Commands & Flags

### Commands:

- 'init': auto build a project (creates packaged artifact) with metadata and logs
- 'config': user defined (inside of a builder.yaml) project build that creates artifact with metadata and logs

### Flags:

- '--help' or '-h': provide info for Builder
- '--output' or '-o': user defined output path for artifact
- '--name' or '-n': user defined project name
- '--yes' or '-y': bypass prompts
- '--branch' or '-b': specify repo branch

## Builder Compatibility

### Languages/Frameworks with default build/install commands:

- Golang (go build)
- Node (npm install)
- Java (maven)
- C# (dotnet)
- Python (pipenv)
- Ruby (bundle install)

## Builder.yaml Parameters

If you are specifying a buildFile within the builder.yaml, you MUST include the projectType.

- projectPath: provide path for project to be built
  - ("C:/Users/Name/Projects", etc)
- projectType: provide language/framework being used
  - (Node, Java, Go, Ruby, Python, C#, Ruby)
- buildTool: provide tool used to install dependencies/build project
  - (maven, npm, bundler, pipenv, etc)
- buildFile: provide file name needed to install dep/build project
  - Can be any user specified file. (myCoolProject.go, package.json etc)
- buildCmd: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)

## Builder ENV Vars

### Native env vars:

- "BUILDER_PARENT_DIR": parent dir path
- "BUILDER_HIDDEN_DIR": hidden dir path
- "BUILDER_LOGS_DIR": logs dir path
- "BUILDER_OUTPUT_PATH": artifact output path

### Envs set by builder.config:

- "BUILDER_DIR_PATH": user defined parent dir path for specific build
- "BUILDER_PROJECT_TYPE": user defined project type (go, java, etc)
- "BUILDER_BUILD_TOOL": user defined build tool (maven, gradle, npm, yarn, etc)
- "BUILDER_BUILD_FILE": user defined build file (myCoolProject.go)

## Builder Signal Flow/Layout

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
  - check for '-y' flag to bypassPrompt
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
  - UserName
  - HomeDir
  - IP
  - Timestamp
  - GitHash

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

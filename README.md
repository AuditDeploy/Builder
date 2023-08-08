<<<<<<< HEAD
# Builder



## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

Already a pro? Just edit this README.md and make it your own. Want to make it easy? [Use the template at the bottom](#editing-this-readme)!

## Add your files

- [ ] [Create](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#create-a-file) or [upload](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#upload-a-file) files
- [ ] [Add files using the command line](https://docs.gitlab.com/ee/gitlab-basics/add-file.html#add-a-file-using-the-command-line) or push an existing Git repository with the following command:

```
cd existing_repo
git remote add origin https://gitlab.com/KatieHarris/builder.git
git branch -M main
git push -uf origin main
```

## Integrate with your tools

- [ ] [Set up project integrations](https://gitlab.com/KatieHarris/builder/-/settings/integrations)

## Collaborate with your team

- [ ] [Invite team members and collaborators](https://docs.gitlab.com/ee/user/project/members/)
- [ ] [Create a new merge request](https://docs.gitlab.com/ee/user/project/merge_requests/creating_merge_requests.html)
- [ ] [Automatically close issues from merge requests](https://docs.gitlab.com/ee/user/project/issues/managing_issues.html#closing-issues-automatically)
- [ ] [Enable merge request approvals](https://docs.gitlab.com/ee/user/project/merge_requests/approvals/)
- [ ] [Set auto-merge](https://docs.gitlab.com/ee/user/project/merge_requests/merge_when_pipeline_succeeds.html)

## Test and Deploy

Use the built-in continuous integration in GitLab.

- [ ] [Get started with GitLab CI/CD](https://docs.gitlab.com/ee/ci/quick_start/index.html)
- [ ] [Analyze your code for known vulnerabilities with Static Application Security Testing(SAST)](https://docs.gitlab.com/ee/user/application_security/sast/)
- [ ] [Deploy to Kubernetes, Amazon EC2, or Amazon ECS using Auto Deploy](https://docs.gitlab.com/ee/topics/autodevops/requirements.html)
- [ ] [Use pull-based deployments for improved Kubernetes management](https://docs.gitlab.com/ee/user/clusters/agent/)
- [ ] [Set up protected environments](https://docs.gitlab.com/ee/ci/environments/protected_environments.html)

***

# Editing this README

When you're ready to make this README your own, just edit this file and use the handy template below (or feel free to structure it however you want - this is just a starting point!). Thank you to [makeareadme.com](https://www.makeareadme.com/) for this template.

## Suggestions for a good README
Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Name
Choose a self-explaining name for your project.

## Description
Let people know what your project can do specifically. Provide context and add a link to any reference visitors might be unfamiliar with. A list of Features or a Background subsection can also be added here. If there are alternatives to your project, this is a good place to list differentiating factors.

## Badges
On some READMEs, you may see small images that convey metadata, such as whether or not all the tests are passing for the project. You can use Shields to add some to your README. Many services also have instructions for adding a badge.

## Visuals
Depending on what you are making, it can be a good idea to include screenshots or even a video (you'll frequently see GIFs rather than actual videos). Tools like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation
Within a particular ecosystem, there may be a common way of installing things, such as using Yarn, NuGet, or Homebrew. However, consider the possibility that whoever is reading your README is a novice and would like more guidance. Listing specific steps helps remove ambiguity and gets people to using your project as quickly as possible. If it only runs in a specific context like a particular programming language version or operating system or has dependencies that have to be installed manually, also add a Requirements subsection.

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## Support
Tell people where they can go to for help. It can be any combination of an issue tracker, a chat room, an email address, etc.

## Roadmap
If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing
State if you are open to contributions and what your requirements are for accepting them.

For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Perhaps there is a script that they should run or some environment variables that they need to set. Make these steps explicit. These instructions could also be useful to your future self.

You can also document commands to lint the code or run tests. These steps help to ensure high code quality and reduce the likelihood that the changes inadvertently break something. Having instructions for running tests is especially helpful if it requires external setup, such as starting a Selenium server for testing in a browser.

## Authors and acknowledgment
Show your appreciation to those who have contributed to the project.

## License
For open source projects, say how it is licensed.

## Project status
If you have run out of energy or time for your project, put a note at the top of the README saying that development has slowed down or stopped completely. Someone may choose to fork your project or volunteer to step in as a maintainer or owner, allowing your project to keep going. You can also make an explicit request for maintainers.
=======
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

### Flags:

- '--help' or '-h': provide info for Builder
- '--output' or '-o': user defined output path for artifact
- '--name' or '-n': user defined project name
- '--yes' or '-y': bypass prompts
- '--branch' or '-b': specify repo branch

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

To use other buildtools, buildcommands, or custome buildfiles you must create builder.yaml and run `config`.

## Builder.yaml Parameters

If you are specifying a buildfile, buildtool, or buildcmd within the builder.yaml, you MUST include the projectType.

At this point in time, please include ALL builder.yaml parameters (all keys must be lowercase), even if they are empty. (This will be addressed in the next update)

- projectpath: provide path for project to be built
  - ("/Users/Name/Projects", etc)
- projecttype: provide language/framework being used
  - (Node, Java, Go, Ruby, Python, C#, Ruby)
- buildtool: provide tool used to install dependencies/build project
  - (maven, npm, bundler, pipenv, etc)
- buildfile: provide file name needed to install dep/build project
  - Can be any user specified file. (myCoolProject.go, package.json etc)
- buildcmd: provide full command to build/compile project
  - ("npm install --silent", "mvn -o package", anything not provided by the Builder as a default)
- outputpath: provide path for artifact to be sent
  - ("/Users/Name/Artifacts", etc)
- globallogs: specify path to global logs
  - ("var/logs/global-logs/logs.txt")
- dockercmd: specify docker command, if building a container
  - ("docker build -t my-project:1.3 .")
- repoBranch: specify repo branch name
  - (“feature/“new-branch”)
- bypassPrompts: bypass prompts
  - (true)

## Builder ENV Vars

### Native env vars:

- "BUILDER_PARENT_DIR": parent dir path
- "BUILDER_HIDDEN_DIR": hidden dir path
- "BUILDER_LOGS_DIR": logs dir path
- "BUILDER_COMMAND": bool if builder cmd is running

### Envs set by builder.config:

- "BUILDER_DIR_PATH": user defined parent dir path for specific build
- "BUILDER_PROJECT_TYPE": user defined project type (go, java, etc)
- "BUILDER_BUILD_TOOL": user defined build tool (maven, gradle, npm, yarn, etc)
- "BUILDER_BUILD_FILE": user defined build file (myCoolProject.go)
- "BUILDER_BUILD_COMMAND": user defined build commmand (yarn install)
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
>>>>>>> 96d867683aba6607e00e55ac6c973925b66f2a88

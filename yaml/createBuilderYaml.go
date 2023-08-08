package yaml

import (
	"Builder/utils/log"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectName   string
	ProjectPath   string
	ProjectType   string
	BuildTool     string
	BuildFile     string
	BuildCmd      string
	OutputPath    string
	GlobalLogs    string
	DockerCmd     string
	RepoBranch    string
	BypassPrompts string
}

func CreateBuilderYaml(fullPath string) {

	projectName := os.Getenv("BUILDER_DIR_NAME")
	projectPath := os.Getenv("BUILDER_DIR_PATH")
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	buildTool := os.Getenv("BUILDER_BUILD_TOOL")
	buildFile := os.Getenv("BUILDER_BUILD_FILE")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
	globalLogs := os.Getenv("GLOBAL_LOGS_PATH")
	dockerCmd := os.Getenv("BUILDER_DOCKER_CMD")
	repoBranch := os.Getenv("REPO_BRANCH")
	bypassPrompts := os.Getenv("BYPASS_PROMPTS")

	builderData := BuilderYaml{
		ProjectName:   projectName,
		ProjectPath:   projectPath,
		ProjectType:   projectType,
		BuildTool:     buildTool,
		BuildFile:     buildFile,
		BuildCmd:      buildCmd,
		OutputPath:    outputPath,
		GlobalLogs:    globalLogs,
		DockerCmd:     dockerCmd,
		RepoBranch:    repoBranch,
		BypassPrompts: bypassPrompts,
	}

	_, err := os.Stat(fullPath + "/builder.yaml")
	if err == nil {
		log.Warn("builder.yaml already exists ⛔️")
	} else {
		OutputData(fullPath, &builderData)
		log.Info("builder.yaml created ✅")
	}
}

func OutputData(fullPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)
	err := ioutil.WriteFile(fullPath+"/builder.yaml", yamlData, 0644)

	if err != nil {
		log.Fatal("builder.yaml creation failed ⛔️")
	}
}

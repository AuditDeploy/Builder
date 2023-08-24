package yaml

import (
	"Builder/spinner"
	"os"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectName   string
	ProjectPath   string
	ProjectType   string
	BuildTool     string
	BuildFile     string
	PreBuildCmd   string
	ConfigCmd     string
	BuildCmd      string
	ArtifactList  string
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
	preBuildCmd := os.Getenv("BUILDER_PREBUILD_COMMAND")
	configCmd := os.Getenv("BUILDER_CONFIG_COMMAND")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	artifactList := os.Getenv("BUILDER_ARTIFACT_LIST")
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
		PreBuildCmd:   preBuildCmd,
		ConfigCmd:     configCmd,
		BuildCmd:      buildCmd,
		ArtifactList:  artifactList,
		OutputPath:    outputPath,
		GlobalLogs:    globalLogs,
		DockerCmd:     dockerCmd,
		RepoBranch:    repoBranch,
		BypassPrompts: bypassPrompts,
	}

	_, err := os.Stat(fullPath + "/builder.yaml")
	if err != nil {
		OutputData(fullPath, &builderData)
		spinner.LogMessage("builder.yaml created ✅", "info")
	}
}

func OutputData(fullPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)
	err := os.WriteFile(fullPath+"/builder.yaml", yamlData, 0644)

	if err != nil {
		spinner.LogMessage("builder.yaml creation failed ⛔️", "fatal")
	}
}

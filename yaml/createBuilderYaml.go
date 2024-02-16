package yaml

import (
	"Builder/spinner"
	"os"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectName  string
	ProjectPath  string
	ProjectType  string
	BuildsDir    string
	BuildTool    string
	BuildFile    string
	PreBuildCmd  string
	ConfigCmd    string
	BuildCmd     string
	ArtifactList string
	OutputPath   string
	DockerCmd    string
	RepoBranch   string
	Push         map[string]interface{}
	AppIcon      string
}

func CreateBuilderYaml(fullPath string) {
	projectName := os.Getenv("BUILDER_DIR_NAME")
	projectPath := os.Getenv("BUILDER_DIR_PATH")
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	buildsDir := os.Getenv("BUILDER_BUILDS_DIR")
	buildTool := os.Getenv("BUILDER_BUILD_TOOL")
	buildFile := os.Getenv("BUILDER_BUILD_FILE")
	preBuildCmd := os.Getenv("BUILDER_PREBUILD_COMMAND")
	configCmd := os.Getenv("BUILDER_CONFIG_COMMAND")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	artifactList := os.Getenv("BUILDER_ARTIFACT_LIST")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")
	dockerCmd := os.Getenv("BUILDER_DOCKER_CMD")
	repoBranch := os.Getenv("REPO_BRANCH")
	var push map[string]interface{}
	if os.Getenv("BUILDER_PUSH_AUTO") != "" {
		push = map[string]interface{}{
			"url":  os.Getenv("BUILDER_PUSH_URL"),
			"auto": os.Getenv("BUILDER_PUSH_AUTO"),
		}
	} else {
		push = map[string]interface{}{
			"url": os.Getenv("BUILDER_PUSH_URL"),
		}
	}
	appIcon := os.Getenv("BUILD_APP_ICON")

	builderData := BuilderYaml{
		ProjectName:  projectName,
		ProjectPath:  projectPath,
		ProjectType:  projectType,
		BuildsDir:    buildsDir,
		BuildTool:    buildTool,
		BuildFile:    buildFile,
		PreBuildCmd:  preBuildCmd,
		ConfigCmd:    configCmd,
		BuildCmd:     buildCmd,
		ArtifactList: artifactList,
		OutputPath:   outputPath,
		DockerCmd:    dockerCmd,
		RepoBranch:   repoBranch,
		Push:         push,
		AppIcon:      appIcon,
	}

	OutputData(fullPath, &builderData)
}

func OutputData(fullPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)

	err := os.WriteFile(fullPath+"/builder.yaml", yamlData, 0755)
	if err != nil {
		spinner.LogMessage("builder.yaml creation failed ⛔️: "+err.Error(), "fatal")
	}

	spinner.LogMessage("builder.yaml created ✅", "info")
}

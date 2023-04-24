package yaml

import (
	"Builder/logger"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectPath string
	ProjectType string
	BuildTool   string
	BuildFile   string
	BuildCmd    string
	OutputPath  string
}

func CreateBuilderYaml(fullPath string) {

	projectPath := os.Getenv("BUILDER_DIR_PATH")
	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	buildTool := os.Getenv("BUILDER_BUILD_TOOL")
	buildFile := os.Getenv("BUILDER_BUILD_FILE")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")
	outputPath := os.Getenv("BUILDER_OUTPUT_PATH")

	builderData := BuilderYaml{
		ProjectPath: projectPath,
		ProjectType: projectType,
		BuildTool:   buildTool,
		BuildFile:   buildFile,
		BuildCmd:    buildCmd,
		OutputPath:  outputPath,
	}

	_, err := os.Stat(fullPath + "/builder.yaml")
	if err == nil {
		logger.WarningLogger.Println("builder.yaml already exists")
	} else {
		OutputData(fullPath, &builderData)
		logger.WarningLogger.Println("Default builder.yaml created")
	}
}

func OutputData(fullPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)
	err := ioutil.WriteFile(fullPath+"/builder.yaml", yamlData, 0644)

	if err != nil {
		logger.ErrorLogger.Println("builder.yaml creation unsuccessful.")
		panic(err)
	}
}

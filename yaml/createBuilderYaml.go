package yaml

import (
	"builder/logger"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type builderYaml struct {
	ProjectPath string
	ProjectType string
	BuildTool   string
	BuildFile   string
	BuildCmd    string
	OutputPath  string
}

func CreatebuilderYaml(fullPath string) {

	projectPath := os.Getenv("builder_DIR_PATH")
	projectType := os.Getenv("builder_PROJECT_TYPE")
	buildTool := os.Getenv("builder_BUILD_TOOL")
	buildFile := os.Getenv("builder_BUILD_FILE")
	buildCmd := os.Getenv("builder_BUILD_COMMAND")
	outputPath := os.Getenv("builder_OUTPUT_PATH")

	builderData := builderYaml{
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

func OutputData(fullPath string, allData *builderYaml) {
	yamlData, _ := yaml.Marshal(allData)
	err := ioutil.WriteFile(fullPath+"/builder.yaml", yamlData, 0644)

	if err != nil {
		logger.ErrorLogger.Println("builder.yaml creation unsuccessful.")
		panic(err)
	}
}

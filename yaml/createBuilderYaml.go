package yaml

import (
	"Builder/logger"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectType string 
	BuildTool 	string
	BuildFile  	string
	BuildCmd		string
}

func CreateBuilderYaml(projectPath string) {

	projectType := os.Getenv("BUILDER_PROJECT_TYPE")
	buildTool := os.Getenv("BUILDER_BUILD_TOOL")
	buildFile := os.Getenv("BUILDER_BUILD_FILE")
	buildCmd := os.Getenv("BUILDER_BUILD_COMMAND")

	builderData := BuilderYaml {
		ProjectType: 	projectType, 
		BuildTool: 		buildTool,
		BuildFile: 		buildFile,
		BuildCmd: 		buildCmd}

	_, err := os.Stat(projectPath+"/builder.yaml")
	if err == nil {
		logger.WarningLogger.Println("builder.yaml already exists")
	} else {
		OutputData(projectPath, &builderData)
		logger.WarningLogger.Println("Default builder.yaml created")
	}
}

func OutputData(projectPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)
	fmt.Println(allData)
	err := ioutil.WriteFile(projectPath+"/builder.yaml", yamlData, 0644)

	if err != nil {
		logger.ErrorLogger.Println("builder.yaml creation unsuccessful.")
		panic(err)
	}
}
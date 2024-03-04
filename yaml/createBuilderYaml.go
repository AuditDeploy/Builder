package yaml

import (
	"Builder/spinner"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type BuilderYaml struct {
	ProjectName           string
	ProjectPath           string
	ProjectType           string
	BuildsDir             string
	BuildTool             string
	BuildFile             string
	PreBuildCmd           string
	ConfigCmd             string
	BuildCmd              string
	ArtifactList          string
	OutputPath            string
	RepoBranch            string
	Docker                map[string]interface{}
	Push                  map[string]interface{}
	AppIcon               string
	ContainerPort         int
	ServicePort           int
	CandidateDependencies []string
	ReleaseEnvs           []EnvData
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
	repoBranch := os.Getenv("REPO_BRANCH")
	var docker map[string]interface{}
	if os.Getenv("BUILDER_DOCKERFILE") != "" && os.Getenv("BUILDER_DOCKER_REGISTRY") != "" && os.Getenv("BUILDER_DOCKER_VERSION") != "" {
		docker = map[string]interface{}{
			"dockerfile": os.Getenv("BUILDER_DOCKERFILE"),
			"registry":   os.Getenv("BUILDER_DOCKER_REGISTRY"),
			"version":    os.Getenv("BUILDER_DOCKER_VERSION"),
		}
	} else if os.Getenv("BUILDER_DOCKERFILE") != "" && os.Getenv("BUILDER_DOCKER_REGISTRY") != "" {
		docker = map[string]interface{}{
			"dockerfile": os.Getenv("BUILDER_DOCKERFILE"),
			"registry":   os.Getenv("BUILDER_DOCKER_REGISTRY"),
		}
	} else if os.Getenv("BUILDER_DOCKERFILE") != "" && os.Getenv("BUILDER_DOCKER_VERSION") != "" {
		docker = map[string]interface{}{
			"dockerfile": os.Getenv("BUILDER_DOCKERFILE"),
			"version":    os.Getenv("BUILDER_DOCKER_VERSION"),
		}
	} else if os.Getenv("BUILDER_DOCKER_REGISTRY") != "" && os.Getenv("BUILDER_DOCKER_VERSION") != "" {
		docker = map[string]interface{}{
			"registry": os.Getenv("BUILDER_DOCKER_REGISTRY"),
			"version":  os.Getenv("BUILDER_DOCKER_VERSION"),
		}
	} else if os.Getenv("BUILDER_DOCKERFILE") != "" {
		docker = map[string]interface{}{
			"dockerfile": os.Getenv("BUILDER_DOCKERFILE"),
		}
	} else if os.Getenv("BUILDER_DOCKER_REGISTRY") != "" {
		docker = map[string]interface{}{
			"registry": os.Getenv("BUILDER_DOCKER_REGISTRY"),
		}
	} else if os.Getenv("BUILDER_DOCKER_VERSION") != "" {
		docker = map[string]interface{}{
			"version": os.Getenv("BUILDER_DOCKER_VERSION"),
		}
	} else {
		docker = map[string]interface{}{}
	}

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

	containerPort, _ := strconv.Atoi(os.Getenv("RELEASE_CONTAINER_PORT"))
	servicePort, _ := strconv.Atoi(os.Getenv("RELEASE_SERVICE_PORT"))

	var dependsOnCandidates []string
	if os.Getenv("RELEASE_DEPENDENCIES") == "" {
		dependsOnCandidates = nil
	} else {
		dependsOnCandidates = strings.Split(os.Getenv("RELEASE_DEPENDENCIES"), ",")
	}

	var releaseEnvs []EnvData
	if os.Getenv("RELEASE_ENVS") != "" {
		envPairs := strings.Split(os.Getenv("RELEASE_ENVS"), ";")
		for _, pair := range envPairs {
			pairArray := strings.Split(pair, ",")
			pairData := EnvData{
				Key:   pairArray[0],
				Value: pairArray[1],
			}
			releaseEnvs = append(releaseEnvs, pairData)
		}
	} else {
		releaseEnvs = nil
	}

	builderData := BuilderYaml{
		ProjectName:           projectName,
		ProjectPath:           projectPath,
		ProjectType:           projectType,
		BuildsDir:             buildsDir,
		BuildTool:             buildTool,
		BuildFile:             buildFile,
		PreBuildCmd:           preBuildCmd,
		ConfigCmd:             configCmd,
		BuildCmd:              buildCmd,
		ArtifactList:          artifactList,
		OutputPath:            outputPath,
		RepoBranch:            repoBranch,
		Docker:                docker,
		Push:                  push,
		AppIcon:               appIcon,
		ContainerPort:         containerPort,
		ServicePort:           servicePort,
		CandidateDependencies: dependsOnCandidates,
		ReleaseEnvs:           releaseEnvs,
	}

	OutputData(fullPath, &builderData)
}

func UpdateBuilderYaml(fullPath string) {

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
	repoBranch := os.Getenv("REPO_BRANCH")
	docker := map[string]interface{}{
		"dockerfile": os.Getenv("BUILDER_DOCKERFILE"),
		"registry":   os.Getenv("BUILDER_DOCKER_REGISTRY"),
		"version":    os.Getenv("BUILDER_DOCKER_VERSION"),
	}

	var push map[string]interface{}
	if os.Getenv("BUILDER_PUSH_AUTO") == "" {
		push = map[string]interface{}{
			"url":  os.Getenv("BUILDER_PUSH_URL"),
			"auto": "false",
		}
	} else {
		push = map[string]interface{}{
			"url":  os.Getenv("BUILDER_PUSH_URL"),
			"auto": os.Getenv("BUILDER_PUSH_AUTO"),
		}
	}
	appIcon := os.Getenv("BUILD_APP_ICON")

	containerPort, _ := strconv.Atoi(os.Getenv("RELEASE_CONTAINER_PORT"))
	servicePort, _ := strconv.Atoi(os.Getenv("RELEASE_SERVICE_PORT"))

	var dependsOnCandidates []string
	if os.Getenv("RELEASE_DEPENDENCIES") == "" {
		dependsOnCandidates = nil
	} else {
		dependsOnCandidates = strings.Split(os.Getenv("RELEASE_DEPENDENCIES"), ",")
	}

	var releaseEnvs []EnvData
	if os.Getenv("RELEASE_ENVS") != "" {
		envPairs := strings.Split(os.Getenv("RELEASE_ENVS"), ";")
		for _, pair := range envPairs {
			pairArray := strings.Split(pair, ",")
			pairData := EnvData{
				Key:   pairArray[0],
				Value: pairArray[1],
			}
			releaseEnvs = append(releaseEnvs, pairData)
		}
	} else {
		releaseEnvs = nil
	}

	builderData := BuilderYaml{
		ProjectName:           projectName,
		ProjectPath:           projectPath,
		ProjectType:           projectType,
		BuildsDir:             buildsDir,
		BuildTool:             buildTool,
		BuildFile:             buildFile,
		PreBuildCmd:           preBuildCmd,
		ConfigCmd:             configCmd,
		BuildCmd:              buildCmd,
		ArtifactList:          artifactList,
		OutputPath:            outputPath,
		RepoBranch:            repoBranch,
		Docker:                docker,
		Push:                  push,
		AppIcon:               appIcon,
		ContainerPort:         containerPort,
		ServicePort:           servicePort,
		CandidateDependencies: dependsOnCandidates,
		ReleaseEnvs:           releaseEnvs,
	}

	_, err := os.Stat(fullPath + "/builder.yaml")
	if err == nil {
		OutputData(fullPath, &builderData)
		spinner.LogMessage("builder.yaml updated ✅", "info")
	}
}

func OutputData(fullPath string, allData *BuilderYaml) {
	yamlData, _ := yaml.Marshal(allData)

	err := os.WriteFile(fullPath+"/builder.yaml", yamlData, 0755)
	if err != nil {
		spinner.LogMessage("builder.yaml creation failed ⛔️: "+err.Error(), "fatal")
	}

	spinner.LogMessage("builder.yaml created ✅", "info")
}

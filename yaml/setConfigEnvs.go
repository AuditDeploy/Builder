package yaml

import (
	"Builder/utils"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type EnvData struct {
	Key   string
	Value string
}

func ConfigEnvs(byi interface{}) {
	//change interface{} into map interface{}
	bldyml, _ := byi.(map[string]interface{})

	//~~~Check for specific key and create env var based on value~~~

	//check for dir name
	if val, ok := bldyml["projectname"]; ok {
		_, present := os.LookupEnv("BUILDER_DIR_NAME")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_DIR_NAME", valStr)
		}
	} else {
		os.Setenv("BUILDER_DIR_NAME", "")
	}

	//check for dir path
	if val, ok := bldyml["projectpath"]; ok {
		_, present := os.LookupEnv("BUILDER_DIR_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)

			// If on windows and they specify a path that begins with '/' append to home dir
			if runtime.GOOS == "windows" && strings.HasPrefix(valStr, "/") {
				homeDir := utils.GetUserData().HomeDir
				os.Setenv("BUILDER_DIR_PATH", homeDir+valStr)
			} else {
				os.Setenv("BUILDER_DIR_PATH", valStr)
			}
		}
	} else {
		os.Setenv("BUILDER_DIR_PATH", "")
	}

	//check for project type
	if val, ok := bldyml["projecttype"]; ok {
		_, present := os.LookupEnv("BUILDER_PROJECT_TYPE")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_PROJECT_TYPE", valStr)
		}
	} else {
		os.Setenv("BUILDER_PROJECT_TYPE", "")
	}

	//check for different dir name to store builds
	if val, ok := bldyml["buildsdir"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILDS_DIR")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILDS_DIR", valStr)
		}
	} else {
		os.Setenv("BUILDER_BUILDS_DIR", "")
	}

	//check for build type
	if val, ok := bldyml["buildtool"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_TOOL")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_TOOL", valStr)
		}
	}

	//check for build file
	if val, ok := bldyml["buildfile"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_FILE")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_FILE", valStr)
		}
	} else {
		os.Setenv("BUILDER_BUILD_FILE", "")
	}

	//check for config cmd
	if val, ok := bldyml["prebuildcmd"]; ok {
		_, present := os.LookupEnv("BUILDER_PREBUILD_COMMAND")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_PREBUILD_COMMAND", valStr)
		}
	}

	//check for config cmd
	if val, ok := bldyml["configcmd"]; ok {
		_, present := os.LookupEnv("BUILDER_CONFIG_COMMAND")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_CONFIG_COMMAND", valStr)
		}
	}

	//check for build cmd
	if val, ok := bldyml["buildcmd"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_COMMAND")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_COMMAND", valStr)
		}
	}

	//check for output path
	if val, ok := bldyml["outputpath"]; ok {
		_, present := os.LookupEnv("BUILDER_OUTPUT_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)

			// If on windows and they specify a path that begins with '/' append to home dir
			if runtime.GOOS == "windows" && strings.HasPrefix(valStr, "/") {
				homeDir := utils.GetUserData().HomeDir
				os.Setenv("BUILDER_OUTPUT_PATH", homeDir+valStr)
			} else {
				os.Setenv("BUILDER_OUTPUT_PATH", valStr)
			}
		}
	}

	//check for docker cmd
	if val, ok := bldyml["dockercmd"]; ok {
		_, present := os.LookupEnv("BUILDER_DOCKER_CMD")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_DOCKER_CMD", valStr)
		}
	}

	//check for global logs path
	if val, ok := bldyml["giturl"]; ok {
		_, present := os.LookupEnv("GIT_URL")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("GIT_URL", valStr)
		}
	}

	//check for an artifacts list
	if val, ok := bldyml["artifactlist"]; ok {
		_, present := os.LookupEnv("BUILDER_ARTIFACT_LIST")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_ARTIFACT_LIST", valStr)
		}
	}

	//check for branch repo
	if val, ok := bldyml["repobranch"]; ok {
		_, present := os.LookupEnv("REPO_BRANCH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("REPO_BRANCH", valStr)
		}
	}

	//check for options to build docker image
	if val, ok := bldyml["docker"]; ok {
		switch v := val.(type) {
		case []interface{}:
			dockerfile := v[0].(map[string]interface{})["dockerfile"]
			registry := v[0].(map[string]interface{})["registry"]
			version := v[0].(map[string]interface{})["version"]

			if dockerfile != nil && dockerfile != "" {
				os.Setenv("BUILDER_DOCKERFILE", dockerfile.(string))
			}
			if registry != nil && registry != "" {
				os.Setenv("BUILDER_DOCKER_REGISTRY", registry.(string))
			}
			if version != nil && version != "" {
				os.Setenv("BUILDER_DOCKER_VERSION", version.(string))
			}
		default: // type map[string]interface{}
			dockerfile := val.(map[string]interface{})["dockerfile"]
			registry := val.(map[string]interface{})["registry"]
			version := val.(map[string]interface{})["version"]

			if dockerfile != nil && dockerfile != "" {
				os.Setenv("BUILDER_DOCKERFILE", dockerfile.(string))
			}
			if registry != nil && registry != "" {
				os.Setenv("BUILDER_DOCKER_REGISTRY", registry.(string))
			}
			if version != nil && version != "" {
				os.Setenv("BUILDER_DOCKER_VERSION", version.(string))
			}
		}
	}

	//check for options to push resulting build data
	if val, ok := bldyml["push"]; ok {
		switch v := val.(type) {
		case []interface{}:
			url := v[0].(map[string]interface{})["url"]
			auto := v[0].(map[string]interface{})["auto"]

			if url != nil && url != "" {
				os.Setenv("BUILDER_PUSH_URL", url.(string))
			}
			if auto != nil {
				os.Setenv("BUILDER_PUSH_AUTO", auto.(string))
			}
		default: // type map[string]interface{}
			url := val.(map[string]interface{})["url"]
			auto := val.(map[string]interface{})["auto"]

			if url != nil && url != "" {
				os.Setenv("BUILDER_PUSH_URL", url.(string))
			}
			if auto != nil {
				os.Setenv("BUILDER_PUSH_AUTO", auto.(string))
			}
		}
	}

	//check for app icon for build
	if val, ok := bldyml["appicon"]; ok {
		_, present := os.LookupEnv("BUILD_APP_ICON")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILD_APP_ICON", valStr)
		}
	}

	// check for container port for application
	if val, ok := bldyml["containerport"]; ok {
		_, present := os.LookupEnv("APP_CONTAINER_PORT")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("APP_CONTAINER_PORT", valStr)
		}
	}

	// check for service port for application
	if val, ok := bldyml["serviceport"]; ok {
		_, present := os.LookupEnv("APP_SERVICE_PORT")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("APP_SERVICE_PORT", valStr)
		}
	}

	// check for application dependencies
	if val, ok := bldyml["application_dependencies"]; ok {
		_, present := os.LookupEnv("APP_DEPENDENCIES")
		if !present {
			//convert list to string
			var dependenciesStr string

			for i, dependency := range val.([]interface{}) {
				if i == len(val.([]interface{}))-1 {
					dependenciesStr += dependency.(string)
				} else {
					dependenciesStr += dependency.(string) + ","
				}
			}

			os.Setenv("APP_DEPENDENCIES", dependenciesStr)
		}
	}

	// check for env vars
	if val, ok := bldyml["application_envs"]; ok {
		_, present := os.LookupEnv("APP_ENVS")
		if !present {
			//convert list of objects to string
			var envsStr string

			for i, envpair := range val.([]interface{}) {
				pair, _ := envpair.(map[string]interface{})
				envsStr += pair["key"].(string) + ","
				envsStr += pair["value"].(string)

				if i != reflect.ValueOf(val).Len()-1 {
					envsStr += ";"
				}
			}

			os.Setenv("APP_ENVS", envsStr)
		}
	}
}

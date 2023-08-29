package yaml

import (
	"Builder/utils"
	"fmt"
	"os"
	"runtime"
	"strings"
)

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
			if runtime.GOOS == "windows" && strings.HasPrefix(valStr, "/") == true {
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
			if runtime.GOOS == "windows" && strings.HasPrefix(valStr, "/") == true {
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
	if val, ok := bldyml["globallogs"]; ok {
		_, present := os.LookupEnv("GLOBAL_LOGS_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("GLOBAL_LOGS_PATH", valStr)
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
}

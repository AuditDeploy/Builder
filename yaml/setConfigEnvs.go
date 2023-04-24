package yaml

import (
	"fmt"
	"os"
)

func ConfigEnvs(byi interface{}) {
	//change interface{} into map interface{}
	bldyml, _ := byi.(map[string]interface{})

	//~~~Check for specific key and create env var based on value~~~

	//check for dir path
	if val, ok := bldyml["projectpath"]; ok {
		_, present := os.LookupEnv("BUILDER_DIR_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_DIR_PATH", valStr)
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

	//check for build file
	if val, ok := bldyml["buildcmd"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_COMMAND")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_COMMAND", valStr)
		}
	}

	//check for build type
	if val, ok := bldyml["outputpath"]; ok {
		_, present := os.LookupEnv("BUILDER_OUTPUT_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_OUTPUT_PATH", valStr)
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
}

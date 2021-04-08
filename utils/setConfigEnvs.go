package utils

import (
	"fmt"
	"os"
)


func ConfigEnvs(byi interface{}) {
	//change interface{} into map interface{}
	bldyml, _ := byi.(map[string]interface{})

	//~~~Check for specific key and create env var based on value~~~
	
	//check for dir path
	if val, ok := bldyml["path"]; ok {
		_, present := os.LookupEnv("BUILDER_DIR_PATH")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_DIR_PATH", valStr)
		}
	}

	//check for project type
	if val, ok := bldyml["projectType"]; ok {
		_, present := os.LookupEnv("BUILDER_PROJECT_TYPE")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_PROJECT_TYPE", valStr)
		} 
	}

	//check for build type
	if val, ok := bldyml["buildTool"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_TOOL")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_TOOL", valStr)
		}
	}

	//check for build file
	if val, ok := bldyml["buildFile"]; ok {
		_, present := os.LookupEnv("BUILDER_BUILD_FILE")
		if !present {
			//convert val interface{} to string to be set as env var
			valStr := fmt.Sprintf("%v", val)
			os.Setenv("BUILDER_BUILD_FILE", valStr)
		}
	}
}
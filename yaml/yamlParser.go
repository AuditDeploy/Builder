package yaml

import (
	"Builder/spinner"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func YamlParser(yamlPath string) {
	// Initial declaration
	m := map[string]interface{}{
		"key": "value",
	}
	// Dynamically add a sub-map
	m["sub"] = map[string]interface{}{
		"deepKey": "deepValue",
	}
	// returns map
	var f interface{}

	//takes yaml path and read file
	source, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		removeTempDir()
		spinner.LogMessage("failed to read builder yaml: "+err.Error(), "fatal")
	}

	//unpacks yaml file in a map int{}
	err = yaml.Unmarshal([]byte(source), &f)
	if err != nil {
		spinner.LogMessage(err.Error(), "error")
	}

	//pass map int{} to callback that sets env vars
	ConfigEnvs(f)

	// if env var BUILDER_COMMAND != true
	removeTempDir()
	//else
}

func removeTempDir() {
	//delete tempRepo dir
	err := os.RemoveAll("./tempRepo")
	if err != nil {
		spinner.LogMessage("Failed to delete tempRepo directory: "+err.Error(), "fatal")
	}
}

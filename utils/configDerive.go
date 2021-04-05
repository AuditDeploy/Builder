package utils

import (
	"os"
	"strings"
)

//ConfigDerive checks "BUILDER_PROJECT_TYPE" env var and returns string arr based on type
func ConfigDerive() []string {
	//make type lowercase
	configType := strings.ToLower(os.Getenv("BUILDER_PROJECT_TYPE"))

	var files []string
	if (configType == "go") {
		files = []string{"main.go"}
	} else if (configType == "node" || configType == "npm") {
		files = []string{"package.json"}
	} else if (configType == "java") {
		files = []string{"pom.xml"}
	} else if (configType == "ruby") {
		files = []string{"gemfile.lock"}
	} else if (configType == "c#" || configType == "csharp") {
		files = []string{".csproj", ".sln"}
	}

	return files
}
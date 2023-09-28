package utils

import (
	"os"
	"strings"
)

// ConfigDerive checks "BUILDER_PROJECT_TYPE" env var and returns string arr based on type
func ConfigDerive() []string {

	//make type lowercase
	configType := strings.ToLower(os.Getenv("BUILDER_PROJECT_TYPE"))
	buildFile := os.Getenv("BUILDER_BUILD_FILE")

	var files []string
	if configType == "go" {
		if buildFile != "" {
			//custom build file from builder.yaml
			files = []string{buildFile}
		} else {
			//default
			files = []string{"main.go"}
		}
	} else if configType == "rust" {
		if buildFile != "" {
			//custom build file from builder.yaml
			files = []string{buildFile}
		} else {
			//default
			files = []string{"Cargo.toml"}
		}
	} else if configType == "node" || configType == "npm" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{"package.json"}
		}
	} else if configType == "java" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{"pom.xml"}
		}
	} else if configType == "ruby" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{"gemfile.lock", "gemfile"}
		}
	} else if configType == "c#" || configType == "csharp" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{".csproj", ".sln"}
		}
	} else if configType == "python" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{"requirements.txt"}
		}
	} else if configType == "c" {
		if buildFile != "" {
			files = []string{buildFile}
		} else {
			files = []string{"Makefile.am"}
		}
	}

	return files
}

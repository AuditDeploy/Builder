package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//CheckArgs is...
func CheckArgs() {
	//Repo
	repo := GetRepoURL()
	cArgs := os.Args[1:]
	//if flag present, but no url
	if repo == "" {
		log.Fatal("No Repo Url Provided")
	}

	//check to see if repo exists
	//git ls-remote lists refs/heads & tags of a repo, if none exists, exit status thrown
	//returns the exit status in err
	_, err := exec.Command("git", "ls-remote", repo, "-q").Output()
	if err != nil {
		log.Fatal("Repo Provided Does Not Exists")
	}

	//check if artifact path is passed in
	var artifactPath string
	for i, v := range cArgs {
		if v == "--output" || v == "-o" {
			if len(cArgs) <= i+1 {
				log.Fatal("No Output Path Provided")

			} else {
				artifactPath = cArgs[i+1]
				val, present := os.LookupEnv("BUILDER_OUTPUT_PATH")
				if !present {
					os.Setenv("BUILDER_OUTPUT_PATH", artifactPath)
				} else {
					fmt.Println("BUILDER_OUTPUT_PATH", val)
					fmt.Println("Output Path already present")
				}
			}
		}

		if strings.Contains(v, "--auto-push=") {
			command := strings.Split(v, "=")[1]
			configType := "json"
			currentDirectory, _ := os.Getwd()
			configDirectory := filepath.Join(currentDirectory, "configs")

			var Cfg Config

			name := GetRepoName(repo)
			filename, v := CreateConfigFile(name+"_config", configType)

			// if v == nil; file already exists; get config from configuration slice created in main
			if v == nil {
				cfgs := RetrieveConfis()

				for _, c := range cfgs {
					if c.ConfigFileUsed() == filename {
						v = c
					}
				}
			}

			allSettings := Cfg.ReadConfigFile(name+"_config.json", configType, configDirectory, v)

			val, ok := allSettings["autopush"]

			if (!ok || val == "0") && command == "true" {
				allSettings["autopush"] = true
				WriteConfigFile(name+"_config.json", configType, configDirectory, &allSettings, v)
				fmt.Println("Successfully updated configuration")
			} else if (!ok || val == "1") && command == "false" {
				allSettings["autopush"] = false
				WriteConfigFile(name+"_config.json", configType, configDirectory, &allSettings, v)
				fmt.Println("Successfully updated configuration")
			}
		}
	}
}

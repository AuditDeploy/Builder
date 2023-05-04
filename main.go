package main

import (
	"Builder/cmd"
	"Builder/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func main() {
	repo := utils.GetRepoURL()
	repoName := utils.GetRepoName(repo)
	currentDirectory, _ := os.Getwd()
	configDirectory := filepath.Join(currentDirectory, "configs")

	utils.InitConfig()

	if len(os.Args) > 1 {
		utils.Help()
		builderCommand := os.Args[1]

		var cfg utils.Config

		if builderCommand == "init" {
			cmd.Init()
			cfgs := utils.RetrieveConfis()
			var config *viper.Viper

			for _, c := range cfgs {
				if c.ConfigFileUsed() == filepath.Join(currentDirectory, "configs", repoName+"_config.json") {
					config = c
				}
			}

			allSettings := cfg.ReadConfigFile(repoName+"_config.json", "json", configDirectory, config)

			push := allSettings["autopush"]

			if push == "1" {
				utils.PushRepo()
			}

		} else if builderCommand == "config" {
			cmd.Config()
			cfgs := utils.RetrieveConfis()
			var config *viper.Viper

			for _, c := range cfgs {
				if c.ConfigFileUsed() == filepath.Join(currentDirectory, "configs", repoName+"_config.json") {
					config = c
				}
			}

			allSettings := cfg.ReadConfigFile(repoName+"_config.json", "json", configDirectory, config)

			push := allSettings["autopush"]

			if push == "1" {
				utils.PushRepo()
			}

		} else {
			cmd.Builder()
		}
	} else {
		cmd.Builder()
	}
	fmt.Println("Build Complete 🔨")
}

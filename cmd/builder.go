package cmd

import (
	"Builder/derive"
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/yaml"
	"net/url"
	"os"
)

func Builder() {
	os.Setenv("BUILDER_COMMAND", "true")
	path, _ := os.Getwd()

	// Start loading spinner
	spinner.Spinner.Start()

	spinner.LogMessage("got here", "info")

	// Check if push command provided and save properties
	args := os.Args
	for i, v := range args {
		if v == "push" {
			if len(args) <= i+1 {
				if os.Getenv("BUILDER_PUSH_URL") == "" {
					spinner.LogMessage("No Push Url Provided", "fatal")
				}
			} else {
				if args[i+1] == "--save" {
					os.Setenv("BUILDER_PUSH_AUTO", "true")

					if len(args) <= i+2 {
						if os.Getenv("BUILDER_PUSH_URL") == "" {
							spinner.LogMessage("No Push Url Provided", "fatal")
						}
					} else {
						pushURL := args[i+2]
						_, err := url.Parse(pushURL)
						if err != nil {
							spinner.LogMessage("Push URL provided is not a valid url: "+err.Error(), "fatal")
						}
						os.Setenv("BUILDER_PUSH_URL", pushURL)
					}
				} else {
					pushURL := args[i+1]
					_, err := url.Parse(pushURL)
					if err != nil {
						spinner.LogMessage("Push URL provided is not a valid url: "+err.Error(), "fatal")
					}
					os.Setenv("BUILDER_PUSH_URL", pushURL)

					if len(args) > i+2 {
						if args[i+2] == "--save" {
							os.Setenv("BUILDER_PUSH_AUTO", "true")
						}
					}
				}
			}
		}
	}

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		// Set repo path
		os.Setenv("BUILDER_REPO_DIR", path)

		// Create directories
		directory.MakeDirs()
		spinner.LogMessage("Directories successfully created.", "info")

		// clone files from current dir into hidden
		currentDir, _ := os.Getwd()
		hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
		utils.CloneRepoFiles(currentDir, hiddenDir)
		spinner.LogMessage("Files copied to hidden dir successfully.", "info")

		//creates a new artifact
		derive.ProjectType()

		//Get build metadata (deprecated, func moved inside compiler)
		spinner.LogMessage("Metadata created successfully.", "info")

		// Store build metadata to hidden builder dir
		utils.StoreBuildMetadataLocally()

		// If provided, push build data to provided url
		if os.Getenv("BUILDER_PUSH_AUTO") == "true" {
			if os.Getenv("BUILDER_PUSH_URL") != "" {
				utils.PushBuildData()
			} else {
				spinner.LogMessage("Build is set to auto-push but was not provided a url.", "fatal")
			}
		} else { // if push property was provided to command, push build data to provided url
			args := os.Args
			for _, v := range args {
				if v == "push" {
					if os.Getenv("BUILDER_PUSH_URL") != "" {
						utils.PushBuildData()
					} else {
						spinner.LogMessage("Push Url Not Provided.", "fatal")
					}
				}
			}
		}

		//Check for Dockerfile, then build image
		utils.Docker()

		//makes hidden dir read-only
		utils.MakeHidden()
		spinner.LogMessage("Hidden Dir is now read-only.", "info")

		// Stop loading spinner
		spinner.Spinner.Stop()
	} else {
		utils.Help()
	}
}

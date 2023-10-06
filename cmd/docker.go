package cmd

import (
	"Builder/spinner"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var closeLocalLogger func()

func Docker() {
	path, _ := os.Getwd()

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger, closeLocalLogger := log.NewLogger("logs", localPath)

	//checks if yaml file exists in path
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		// Start loading spinner
		spinner.Spinner.Start()

		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		name := utils.GetName()
		startTime := time.Now().Unix()
		dockerfile := os.Getenv("BUILDER_DOCKERFILE")
		var cmd *exec.Cmd

		// Docker build cmd
		if dockerfile != "" {
			// user specified Dockerfile
			cmd = exec.Command("docker", "-f", dockerfile, "-t", name+"_"+fmt.Sprint(startTime), ".")
			cmd.Dir = path
		} else {
			//default
			cmd = exec.Command("docker", "build", "-t", name+"_"+fmt.Sprint(startTime), ".")
			cmd.Dir = path
		}

		//run cmd, check for err, log cmd
		spinner.LogMessage("running command: "+cmd.String(), "info")

		stdout, pipeErr := cmd.StdoutPipe()
		if pipeErr != nil {
			spinner.LogMessage(pipeErr.Error(), "fatal")
		}

		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		done := make(chan struct{})

		scanner := bufio.NewScanner(stdout)

		// Use the scanner to scan the output line by line and log it
		// It's running in a goroutine so that it doesn't block
		go func() {
			// Read line by line and process it
			for scanner.Scan() {
				line := scanner.Text()
				// Have to stop spinner or it will get printed with log to console
				spinner.Spinner.Stop()
				locallogger.Info(line)
				spinner.Spinner.Start()
			}

			// We're all done, unblock the channel
			done <- struct{}{}

		}()

		if err := cmd.Start(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Wait for all output to be processed
		<-done

		// Wait for cmd to finish
		if err := cmd.Wait(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// TODO: Get build name from metadata.json file from within container and re tag local image to it
		cmd = exec.Command("docker", "cp", name+"_"+fmt.Sprint(startTime)+":/root/.builder/builds.json", path)
		cmd.Dir = path

		//run cmd, check for err, log cmd
		spinner.LogMessage("retrieving build info from container...", "info")

		if err := cmd.Run(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// If release tag provided push image to user provided remote registry
		args := os.Args

		for _, v := range args {
			if v == "--release" || v == "-r" {
				// TODO: Check for remote registry and tag and push to it
			}
		}

		// Close log file
		closeLocalLogger()

		// Stop loading spinner
		spinner.Spinner.Stop()
	} else {
		utils.Help()
	}
}

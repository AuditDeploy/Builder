package cmd

import (
	"Builder/directory"
	"Builder/spinner"
	"Builder/utils"
	"Builder/utils/log"
	"Builder/yaml"
	"runtime"

	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mongodb/mongo-tools/common/password"
)

var closeLocalLogger func()

type Metadata struct {
	ProjectName       string                 `json:"ProjectName"`
	ProjectType       string                 `json:"ProjectType"`
	ArtifactName      string                 `json:"ArtifactName"`
	ArtifactChecksums map[string]interface{} `json:"ArtifactChecksums"`
	ArtifactLocation  string                 `json:"ArtifactLocation"`
	Logs              string
	LogsLocation      string `json:"LogsLocation"`
	Username          string `json:"UserName"`
	HomeDir           string `json:"HomeDir"`
	IP                string `json:"IP"`
	StartTime         string `json:"StartTime"`
	EndTime           string `json:"EndTime"`
	GitURL            string `json:"GitURL"`
	MasterGitHash     string `json:"MasterGitHash"`
	BranchName        string `json:"BranchName"`
}

func Docker() {
	os.Setenv("BUILDER_DOCKER_COMMAND", "true")
	path, _ := os.Getwd()

	//Set up local logger
	localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
	locallogger, closeLocalLogger := log.NewLogger("logs", localPath)

	// Start loading spinner
	spinner.Spinner.Start()

	//checks if yaml file exists in path, if it does, continue
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		// make dirs
		directory.MakeDirs()
		spinner.LogMessage("Directories successfully created.", "info")

		releaseTagProvided := false
		args := os.Args

		// If user provides release tag, have them docker login
		for _, v := range args {
			if v == "--release" || v == "-r" {
				releaseTagProvided = true

				spinner.LogMessage("Attempting Docker login...", "info")

				if os.Getenv("BUILDER_DOCKER_REGISTRY") == "" {
					spinner.LogMessage("Cannot login to Docker registry: No Docker registry provided, please provide in the builder.yaml", "fatal")
				} else {
					dockerRegistry := os.Getenv("BUILDER_DOCKER_REGISTRY")[:strings.IndexByte(os.Getenv("BUILDER_DOCKER_REGISTRY"), '/')]

					// Try to do a docker login (user might have login saved), if can't login, ask for username and password to login
					var cmd *exec.Cmd
					if runtime.GOOS == "windows" {
						cmd = exec.Command("docker", "login", dockerRegistry)
					} else {
						cmd = exec.Command("sudo", "docker", "login", dockerRegistry)
					}
					if err := cmd.Run(); err != nil {
						// Stop spinner for input from user
						spinner.Spinner.Stop()

						fmt.Println("Please login with your Docker ID or email address to push image to provided registry.")

						// Ask for Docker username
						fmt.Print("Username: ")
						var dockerUsername string
						fmt.Scan(&dockerUsername)

						// Ask for Docker password
						dockerPassword, _ := password.Prompt(dockerUsername)

						// Start loading spinner
						spinner.Spinner.Start()

						// Try to login to provided registry with these login creds
						if err := exec.Command("docker", "login", os.Getenv("BUILDER_DOCKER_REGISTRY"), "--username", dockerUsername, "--password", dockerPassword).Run(); err != nil {
							spinner.LogMessage("Cannot login to Docker registry: "+err.Error(), "fatal")
						}
					}

					spinner.LogMessage("Successfully logged in to Docker registry.", "info")
				}
			}
		}

		name := utils.GetName()
		startTime := time.Now().Unix()
		dockerfile := os.Getenv("BUILDER_DOCKERFILE")
		var cmd *exec.Cmd

		// Docker build cmd
		if dockerfile != "" {
			// user specified Dockerfile
			cmd = exec.Command("docker", "-f", dockerfile, "-t ", name+"_"+fmt.Sprint(startTime), " .")
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

		var copyPath string
		if os.Getenv("BUILDER_DIR_PATH") != "" {
			copyPath = os.Getenv("BUILDER_DIR_PATH")
		} else {
			copyPath = path
		}

		// Create a container from the image to grab build info
		if err := exec.Command("docker", "create", "--name", name+"_"+fmt.Sprint(startTime), name+"_"+fmt.Sprint(startTime)+":latest").Run(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Copy build info from created container
		cmd = exec.Command("docker", "cp", name+"_"+fmt.Sprint(startTime)+":/root/.builder/builds.json", "./builder_builds.json")
		cmd.Dir = copyPath

		//run cmd, check for err, log cmd
		spinner.LogMessage("retrieving build info from container...", "info")

		stdout, pipeErr = cmd.StdoutPipe()
		if pipeErr != nil {
			spinner.LogMessage(pipeErr.Error(), "fatal")
		}

		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		copyMetadataFileFromContainer := make(chan struct{})

		scanner = bufio.NewScanner(stdout)

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
			copyMetadataFileFromContainer <- struct{}{}

		}()

		if err := cmd.Start(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Wait for all output to be processed
		<-copyMetadataFileFromContainer

		// Wait for cmd to finish
		if err := cmd.Wait(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Retrieve build info from file copied from container
		data, err := os.ReadFile(copyPath + "/builder_builds.json")
		if err != nil {
			spinner.LogMessage("Can't open file copied from container: "+err.Error(), "fatal")
		}

		// Decode json into interface
		var metadata Metadata
		json.NewDecoder(strings.NewReader(string(data))).Decode(&metadata)

		// Retrieve logs file from container
		cmd = exec.Command("docker", "cp", name+"_"+fmt.Sprint(startTime)+":"+metadata.LogsLocation, "./builder_logs.json")
		cmd.Dir = copyPath

		//run cmd, check for err, log cmd
		spinner.LogMessage("retrieving build logs from container...", "info")

		stdout, pipeErr = cmd.StdoutPipe()
		if pipeErr != nil {
			spinner.LogMessage(pipeErr.Error(), "fatal")
		}

		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		copyLogsFileFromContainer := make(chan struct{})

		scanner = bufio.NewScanner(stdout)

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
			copyLogsFileFromContainer <- struct{}{}

		}()

		if err := cmd.Start(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Wait for all output to be processed
		<-copyLogsFileFromContainer

		// Wait for cmd to finish
		if err := cmd.Wait(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Retrieve build logs from file copied from container
		logsData, err := os.ReadFile(copyPath + "/builder_logs.json")
		if err != nil {
			spinner.LogMessage("Can't open file copied from container: "+err.Error(), "fatal")
		}

		// Save collected build metadata and logs to file
		SaveBuildMetadata(metadata, string(logsData))

		gatheredProjectName := metadata.ProjectName
		gatheredStartTime, err := time.Parse(time.RFC850, metadata.StartTime)
		if err != nil {
			spinner.LogMessage("Couldn't parse time: "+err.Error(), "fatal")
		}

		// Rename docker container to same name as build completed in container
		if err := exec.Command("docker", "tag", name+"_"+fmt.Sprint(startTime), gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Remove previously tagged image
		if err := exec.Command("docker", "rmi", "-f", name+"_"+fmt.Sprint(startTime)+":latest").Run(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// Remove temp container
		if err := exec.Command("docker", "rm", name+"_"+fmt.Sprint(startTime)).Run(); err != nil {
			spinner.LogMessage(err.Error(), "fatal")
		}

		// If release tag provided push image to user provided remote registry
		if releaseTagProvided == true {
			// Check for remote registry and tag and push to it
			if os.Getenv("BUILDER_DOCKER_REGISTRY") == "" {
				spinner.LogMessage("Cannot complete docker push: No Docker registry provided, please provide in the builder.yaml", "fatal")
			} else {
				spinner.LogMessage("Tagging and pushing docker image...", "info")

				dockerRegistry := os.Getenv("BUILDER_DOCKER_REGISTRY")
				// Re-tag docker image to include remote registry
				if err := exec.Command("docker", "tag", gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix()), dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
				}

				// Remove previously tagged image
				if err := exec.Command("docker", "rmi", "-f", gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())+":latest").Run(); err != nil {
					spinner.LogMessage(err.Error(), "fatal")
				}

				// Push re-tagged docker image to user provided docker registry
				if err := exec.Command("docker", "push", dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage("Could not complete docker push: "+err.Error(), "fatal")
				}
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

func SaveBuildMetadata(metadata Metadata, logsJSON string) {
	// Add logs to metadata
	metadata.Logs = logsJSON

	// Create JSON from metadata struct
	// json, err := json.Marshal(metadata)
	// if err != nil {
	// 	spinner.LogMessage("Could not format build info as json: "+err.Error(), "fatal")
	// }

}

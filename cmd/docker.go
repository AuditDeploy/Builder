package cmd

import (
	"Builder/artifact"
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

	goYaml "gopkg.in/yaml.v2"
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

	// Start loading spinner
	spinner.Spinner.Start()

	// If builder_data folder already exists, remove it
	if _, err := os.Stat(path + "/" + "builder_data"); err == nil {
		e := os.RemoveAll(path + "/" + "builder_data")
		if e != nil {
			spinner.LogMessage("Couldn't remove builder data from previous build: "+e.Error(), "error")
		}
	}

	// make dirs
	directory.MakeDirs()
	spinner.LogMessage("Directories successfully created.", "info")

	//checks if yaml file exists in path, if it does, continue
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		name := utils.GetName()
		startTime := time.Now()
		dockerfile := os.Getenv("BUILDER_DOCKERFILE")
		var cmd *exec.Cmd

		// Docker build cmd
		if dockerfile != "" {
			// user specified Dockerfile
			if runtime.GOOS == "windows" {
				cmd = exec.Command("docker", "-f", dockerfile, "-t", name+"_"+fmt.Sprint(startTime.Unix()), ".")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "sudo docker -f "+dockerfile+" -t "+name+"_"+fmt.Sprint(startTime.Unix())+" .")
			}
			cmd.Dir = path
		} else {
			//default
			if runtime.GOOS == "windows" {
				cmd = exec.Command("docker", "build", "-t", name+"_"+fmt.Sprint(startTime.Unix()), ".")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "sudo docker build -t "+name+"_"+fmt.Sprint(startTime.Unix())+" .")
			}
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

		//Set up local logger
		localPath, _ := os.LookupEnv("BUILDER_LOGS_DIR")
		locallogger, closeLocalLogger := log.NewLogger("docker_logs", localPath)

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

		copyPath := os.Getenv("BUILDER_WORKSPACE_DIR")

		// Create a container from the image to grab build info
		if runtime.GOOS == "windows" {
			if err := exec.Command("docker", "create", "--name", name+"_"+fmt.Sprint(startTime.Unix()), name+"_"+fmt.Sprint(startTime.Unix())+":latest").Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		} else {
			if err := exec.Command("/bin/sh", "-c", "sudo docker create --name "+name+"_"+fmt.Sprint(startTime.Unix())+" "+name+"_"+fmt.Sprint(startTime.Unix())+":latest").Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		}

		// Copy build info from created container
		if runtime.GOOS == "windows" {
			cmd = exec.Command("docker", "cp", name+"_"+fmt.Sprint(startTime.Unix())+":/root/.builder/builds.json", "./builder_builds.json")
		} else {
			cmd = exec.Command("/bin/sh", "-c", "sudo docker cp "+name+"_"+fmt.Sprint(startTime.Unix())+":/root/.builder/builds.json ./builder_builds.json")
		}
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

		gatheredProjectName := metadata.ProjectName
		gatheredStartTime, err := time.Parse(time.RFC850, metadata.StartTime)
		if err != nil {
			spinner.LogMessage("Couldn't parse time: "+err.Error(), "fatal")
		}
		os.Setenv("BUILD_START_TIME", metadata.StartTime)

		// Create dir for metadata
		artifact.ArtifactDir()
		artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

		// Save collected build metadata and docker metadata to file
		SaveBuildMetadata(metadata, artifactDir)

		// Retrieve logs file from container
		if runtime.GOOS == "windows" {
			cmd = exec.Command("docker", "cp", name+"_"+fmt.Sprint(startTime.Unix())+":"+metadata.LogsLocation, "./builder_logs.json")
		} else {
			cmd = exec.Command("/bin/sh", "-c", "sudo docker cp "+name+"_"+fmt.Sprint(startTime.Unix())+":"+metadata.LogsLocation+" ./builder_logs.json")
		}
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
		logsJSON, err := os.ReadFile(copyPath + "/builder_logs.json")
		if err != nil {
			spinner.LogMessage("Can't open file copied from container: "+err.Error(), "fatal")
		}

		logsDir := os.Getenv("BUILDER_LOGS_DIR")

		// Save build logs with docker logs to file
		SaveBuildLogs(logsJSON, logsDir)

		// Close log file
		closeLocalLogger()

		// Rename docker container to same name as build completed in container
		if runtime.GOOS == "windows" {
			if err := exec.Command("docker", "tag", name+"_"+fmt.Sprint(startTime.Unix()), gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		} else {
			if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+name+"_"+fmt.Sprint(startTime.Unix())+" "+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		}

		// Remove previously tagged image
		if runtime.GOOS == "windows" {
			if err := exec.Command("docker", "rmi", "-f", name+"_"+fmt.Sprint(startTime.Unix())+":latest").Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		} else {
			if err := exec.Command("/bin/sh", "-c", "sudo docker rmi -f "+name+"_"+fmt.Sprint(startTime.Unix())+":latest").Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		}

		// Remove temp container
		if runtime.GOOS == "windows" {
			if err := exec.Command("docker", "rm", name+"_"+fmt.Sprint(startTime.Unix())).Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		} else {
			if err := exec.Command("/bin/sh", "-c", "sudo docker rm "+name+"_"+fmt.Sprint(startTime.Unix())).Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		}

		// Update parent directory name
		copyPath = directory.UpdateParentDirName(copyPath)

		// Set start time back to docker start time for when we get docker metadata later
		os.Setenv("BUILD_START_TIME", fmt.Sprint(startTime))

		// If release tag provided push image to user provided remote registry
		args := os.Args
		for _, v := range args {
			if v == "--release" || v == "-r" {
				// Check for remote registry and tag and push to it
				if os.Getenv("BUILDER_DOCKER_REGISTRY") == "" {
					spinner.LogMessage("Cannot complete docker push: No Docker registry provided, please provide in the builder.yaml", "fatal")
				} else {
					spinner.LogMessage("Tagging and pushing docker image...", "info")

					dockerRegistry := os.Getenv("BUILDER_DOCKER_REGISTRY")
					// Re-tag docker image to include remote registry
					if runtime.GOOS == "windows" {
						if err := exec.Command("docker", "tag", gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix()), dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
							spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
						}
					} else {
						if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())+" "+dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
							spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
						}
					}

					// Remove previously tagged image
					if runtime.GOOS == "windows" {
						if err := exec.Command("docker", "rmi", "-f", gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())+":latest").Run(); err != nil {
							spinner.LogMessage(err.Error(), "fatal")
						}
					} else {
						if err := exec.Command("/bin/sh", "-c", "sudo docker rmi -f "+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())+":latest").Run(); err != nil {
							spinner.LogMessage(err.Error(), "fatal")
						}
					}

					// Push re-tagged docker image to user provided docker registry
					if runtime.GOOS == "windows" {
						if err := exec.Command("docker", "push", dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
							spinner.LogMessage("Could not complete docker push: "+err.Error()+".  You may need to docker login.", "fatal")
						}
					} else {
						if err := exec.Command("/bin/sh", "-c", "sudo docker push "+dockerRegistry+"/"+gatheredProjectName+"_"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
							spinner.LogMessage("Could not complete docker push: "+err.Error()+".  You may need to docker login.", "fatal")
						}
					}

					spinner.LogMessage("Docker image successfully tagged and pushed to provided registry.", "info")
				}
			}
		}

		os.Setenv("BUILD_END_TIME", time.Now().String())

		// Create metadata for docker build
		artifactDir = os.Getenv("BUILDER_ARTIFACT_DIR")
		utils.Metadata(artifactDir)
		spinner.LogMessage("Metadata saved successfully.", "info")

		// Stop loading spinner
		spinner.Spinner.Stop()
	} else {
		utils.Help()
	}
}

func SaveBuildMetadata(metadata Metadata, path string) {
	yamlData, _ := goYaml.Marshal(metadata)
	jsonData, _ := json.Marshal(metadata)

	err := os.WriteFile(path+"/metadata.json", jsonData, 0666)
	err2 := os.WriteFile(path+"/metadata.yaml", yamlData, 0666)

	if err != nil {
		spinner.LogMessage("JSON Metadata creation unsuccessful.", "fatal")
	}

	if err2 != nil {
		spinner.LogMessage("YAML Metadata creation unsuccessful.", "fatal")
	}

	// Delete metadata file grabbed from container to workspace dir
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	e := os.Remove(workspaceDir + "/builder_builds.json")
	if e != nil {
		spinner.LogMessage("Couldn't delete temporary log file in workspace: "+e.Error(), "error")
	}
}

func SaveBuildLogs(logsJSON []byte, path string) {
	dockerLogs, err := os.ReadFile(path + "/docker_logs.json")

	f, err := os.OpenFile(path+"/logs.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		spinner.LogMessage("Cannot create logs.json file: "+err.Error(), "fatal")
	}

	if _, err = f.WriteString(string(logsJSON)); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString("\n========== Build End =========="); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString("\n========== Docker Start ==========\n"); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString(string(dockerLogs)); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}

	f.Close()

	// Delete build logs file grabbed from container to workspace dir
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	e := os.Remove(workspaceDir + "/builder_logs.json")
	if e != nil {
		spinner.LogMessage("Couldn't delete temporary log file in workspace: "+e.Error(), "error")
	}
}

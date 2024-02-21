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

// BuildMetadata holds the data gathered from the build done in the docker container
type BuildMetadata struct {
	ProjectName       string `json:"ProjectName"`
	ProjectType       string `json:"ProjectType"`
	ArtifactName      string `json:"ArtifactName"`
	ArtifactChecksums string `json:"ArtifactChecksums"`
	ArtifactLocation  string `json:"ArtifactLocation"`
	LogsLocation      string `json:"LogsLocation"`
	UserName          string `json:"UserName"`
	HomeDir           string `json:"HomeDir"`
	IP                string `json:"IP"`
	StartTime         string `json:"StartTime"`
	EndTime           string `json:"EndTime"`
	GitURL            string `json:"GitURL"`
	MasterGitHash     string `json:"MasterGitHash"`
	BranchName        string `json:"BranchName"`
	BuildID           string `json:"BuildID"`
}

// DockerMetadata holds the data gathered during the docker image build
type DockerMetadata struct {
	ProjectName       string `json:"ProjectName"`
	ProjectType       string `json:"ProjectType"`
	ArtifactName      string `json:"ArtifactName"`
	ArtifactChecksums string `json:"ArtifactChecksums"`
	ArtifactLocation  string `json:"ArtifactLocation"`
	LogsLocation      string `json:"LogsLocation"`
	UserName          string `json:"UserName"`
	HomeDir           string `json:"HomeDir"`
	IP                string `json:"IP"`
	StartTime         string `json:"StartTime"`
	EndTime           string `json:"EndTime"`
	GitURL            string `json:"GitURL"`
	MasterGitHash     string `json:"MasterGitHash"`
	BranchName        string `json:"BranchName"`
}

// AllDockerMetaData holds the data for docker metadata.json file
type AllDockerMetadata struct {
	DockerBuild         BuildMetadata
	DockerImageTag      DockerMetadata
	ProjectName         string
	ProjectType         string
	UserName            string
	StartTime           string
	EndTime             string
	GitURL              string
	MasterGitHash       string
	BuildID             string
	DockerRepositoryTag string
	DockerRepository    string
	DockerTags          []string
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

	//checks if yaml file exists in path, if it does, continue
	if _, err := os.Stat(path + "/" + "builder.yaml"); err == nil {
		//parse builder.yaml
		yaml.YamlParser(path + "/" + "builder.yaml")

		// If --release or -r tag and registry provided, update env var for docker registry
		args := os.Args
		for i, v := range args {
			if v == "--release" || v == "-r" {
				if len(args) <= i+1 {
					spinner.LogMessage("No Docker registry provided.  Please provide to command or in builder.yaml", "fatal")
				} else {
					os.Setenv("BUILDER_DOCKER_REGISTRY", args[i+1])
				}
			}
		}

		// make dirs
		directory.MakeDirs()
		spinner.LogMessage("Directories successfully created.", "info")

		// Start building docker image
		name := os.Getenv("BUILDER_DIR_NAME")
		startTime := time.Now()
		dockerfile := os.Getenv("BUILDER_DOCKERFILE")
		var cmd *exec.Cmd

		// Docker build cmd
		if dockerfile != "" {
			// user specified Dockerfile
			if runtime.GOOS == "windows" {
				cmd = exec.Command("docker", "build", "-f", dockerfile, "-t", name+"_"+fmt.Sprint(startTime.Unix()), ".")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "sudo docker build -f "+dockerfile+" -t "+name+"_"+fmt.Sprint(startTime.Unix())+" .")
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
		var metadata BuildMetadata
		json.NewDecoder(strings.NewReader(string(data))).Decode(&metadata)

		gatheredStartTime, err := time.Parse(time.RFC850, metadata.StartTime)
		if err != nil {
			spinner.LogMessage("Couldn't parse time: "+err.Error(), "fatal")
		}
		os.Setenv("BUILD_START_TIME", metadata.StartTime)
		os.Setenv("BUILD_END_TIME", metadata.EndTime)

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

		// Create dir for metadata
		artifact.ArtifactDir()

		// Close log file
		closeLocalLogger()

		// Rename docker container to same name as build completed in container
		if runtime.GOOS == "windows" {
			if err := exec.Command("docker", "tag", name+"_"+fmt.Sprint(startTime.Unix())+":latest", metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
				spinner.LogMessage(err.Error(), "fatal")
			}
		} else {
			if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+name+"_"+fmt.Sprint(startTime.Unix())+":latest "+metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
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

		// Get additional tags for docker image
		runningTag := fmt.Sprint(gatheredStartTime.Unix()) + "_" + metadata.BuildID + "_" + metadata.UserName
		tags := []string{
			"latest",
			GetHumanReadableStartTimeTag(gatheredStartTime),
			fmt.Sprint(gatheredStartTime.Unix()) + "_" + metadata.BuildID,
			runningTag,
		}
		_, masterGitHash := utils.GitMasterNameAndHash()
		// If masterGitHash exists, add it to a tag
		if masterGitHash != "" && masterGitHash != "undefined" {
			runningTag = runningTag + "_" + masterGitHash[0:7]
			tags = append(tags, runningTag)
		}
		// If user provides version tag, add it as a tag and append it to the running tag
		version := os.Getenv("BUILDER_DOCKER_VERSION")
		if version != "" {
			runningTag = runningTag + "_" + version
			tags = append(tags, version, runningTag)
		}

		// Update parent directory name
		copyPath = directory.UpdateParentDirName(copyPath)

		// If Docker registry provided in the builder.yaml tag and push the image to it
		dockerRegistryProvided := false
		var remoteDockerRepo string
		if os.Getenv("BUILDER_DOCKER_REGISTRY") != "" {
			dockerRegistryProvided = true
			spinner.LogMessage("Tagging and pushing docker image...", "info")

			dockerRegistry := os.Getenv("BUILDER_DOCKER_REGISTRY")
			remoteDockerRepo = dockerRegistry + "/" + metadata.ProjectName
			// Re-tag docker image to include remote registry
			if runtime.GOOS == "windows" {
				if err := exec.Command("docker", "tag", metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix()), remoteDockerRepo+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
				}
			} else {
				if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())+" "+remoteDockerRepo+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
				}
			}

			// Remove previously tagged image
			if runtime.GOOS == "windows" {
				if err := exec.Command("docker", "rmi", "-f", metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage(err.Error(), "fatal")
				}
			} else {
				if err := exec.Command("/bin/sh", "-c", "sudo docker rmi -f "+metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())).Run(); err != nil {
					spinner.LogMessage(err.Error(), "fatal")
				}
			}

			// Add more tags to new remote image
			for _, tag := range tags {
				if runtime.GOOS == "windows" {
					if err := exec.Command("docker", "tag", remoteDockerRepo+":"+fmt.Sprint(gatheredStartTime.Unix()), remoteDockerRepo+":"+tag).Run(); err != nil {
						spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
					}
				} else {
					if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+remoteDockerRepo+":"+fmt.Sprint(gatheredStartTime.Unix())+" "+remoteDockerRepo+":"+tag).Run(); err != nil {
						spinner.LogMessage("Could not re-tag docker image to include registry: "+err.Error(), "fatal")
					}
				}
			}

			// Push re-tagged docker image (and all of its tags) to user provided docker registry
			if runtime.GOOS == "windows" {
				if err := exec.Command("docker", "push", remoteDockerRepo, "--all-tags").Run(); err != nil {
					spinner.LogMessage("Could not complete docker push: "+err.Error()+".  You may need to docker login.", "fatal")
				}
			} else {
				if err := exec.Command("/bin/sh", "-c", "sudo docker push "+remoteDockerRepo+" --all-tags").Run(); err != nil {
					spinner.LogMessage("Could not complete docker push: "+err.Error()+".  You may need to docker login.", "fatal")
				}
			}

			spinner.LogMessage("Docker image successfully tagged and pushed to provided registry.", "info")
		} else { // If not using builder docker release, add the extra tags to local docker image:
			for _, tag := range tags {
				if runtime.GOOS == "windows" {
					if err := exec.Command("docker", "tag", metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix()), metadata.ProjectName+":"+tag).Run(); err != nil {
						spinner.LogMessage("Could not re-tag docker image to include additional tag: "+err.Error(), "fatal")
					}
				} else {
					if err := exec.Command("/bin/sh", "-c", "sudo docker tag "+metadata.ProjectName+":"+fmt.Sprint(gatheredStartTime.Unix())+" "+metadata.ProjectName+":"+tag).Run(); err != nil {
						spinner.LogMessage("Could not re-tag docker image to include additional tag: "+err.Error(), "fatal")
					}
				}
			}
		}

		// Save docker image with running tag and list of tags to env vars for metadata
		if dockerRegistryProvided {
			os.Setenv("BUILDER_DOCKER_REPO", remoteDockerRepo)
			os.Setenv("BUILDER_DOCKER_REPO_TAG", remoteDockerRepo+":"+runningTag)
		} else {
			os.Setenv("BUILDER_DOCKER_REPO", metadata.ProjectName)
			os.Setenv("BUILDER_DOCKER_REPO_TAG", metadata.ProjectName+":"+runningTag)
		}
		os.Setenv("BUILDER_DOCKER_TAGS", fmt.Sprintf("%+v", tags))

		// Create metadata for docker image build
		artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
		OutputDockerMetadata(metadata, tags, artifactDir)
		spinner.LogMessage("Metadata saved successfully.", "info")

		// Remove unneeded docker_logs.json file
		logsDir = os.Getenv("BUILDER_LOGS_DIR")
		e := os.Remove(logsDir + "/docker_logs.json")
		if e != nil {
			spinner.LogMessage("Couldn't delete old docker_logs.json file: "+e.Error(), "error")
		}

		// Re-create builder.yaml to include any new vars
		yaml.UpdateBuilderYaml(path)

		// Stop loading spinner
		spinner.Spinner.Stop()
	} else {
		utils.Help()
	}
}

func GetHumanReadableStartTimeTag(startTime time.Time) string {
	hours := fmt.Sprint(startTime.Hour())
	minutesAsInt := startTime.Minute()
	formattedTime := startTime.Format(time.RFC822)

	// Reformat minutes if 0-9 change to 00-09
	var minutes string
	if minutesAsInt <= 9 {
		stringMinutes := fmt.Sprint(minutesAsInt)
		minutes = fmt.Sprint(0) + stringMinutes
	} else {
		minutes = fmt.Sprint(minutesAsInt)
	}

	// Replace __:__ time in formattedTime with __H__M
	updatedTime := strings.ReplaceAll(formattedTime, hours+":"+minutes, hours+"H"+minutes+"M")

	finalTag := strings.ReplaceAll(updatedTime, " ", "_")

	return finalTag
}

func OutputDockerMetadata(buildMetadata BuildMetadata, imageTags []string, path string) {
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	utils.Metadata(workspaceDir) // Round up all metadata around building docker image

	// Retrieve docker metadata from file created by utils.Metadata()
	data, err := os.ReadFile(workspaceDir + "/metadata.json")
	if err != nil {
		spinner.LogMessage("Can't open metadata file from docker image build: "+err.Error(), "fatal")
	}

	// Decode json into interface
	var dockerMetadata DockerMetadata
	json.NewDecoder(strings.NewReader(string(data))).Decode(&dockerMetadata)

	// Get master git hash of current project
	_, masterGitHash := utils.GitMasterNameAndHash()

	// Retrieve logs

	// Store all docker metadata we want to output to metadata.json
	allDockerMetadata := AllDockerMetadata{
		DockerBuild:         buildMetadata,
		DockerImageTag:      dockerMetadata,
		ProjectName:         buildMetadata.ProjectName,
		ProjectType:         buildMetadata.ProjectType,
		UserName:            utils.GetUserData().Username,
		StartTime:           buildMetadata.StartTime,
		EndTime:             buildMetadata.EndTime,
		GitURL:              utils.GetRepoURL(),
		MasterGitHash:       masterGitHash,
		BuildID:             buildMetadata.BuildID,
		DockerRepositoryTag: os.Getenv("BUILDER_DOCKER_REPO_TAG"),
		DockerRepository:    os.Getenv("BUILDER_DOCKER_REPO"),
		DockerTags:          imageTags,
	}

	jsonData, _ := json.Marshal(allDockerMetadata)
	yamlData, _ := goYaml.Marshal(allDockerMetadata)

	jsonErr := os.WriteFile(path+"/metadata.json", jsonData, 0666)
	yamlErr := os.WriteFile(path+"/metadata.yaml", yamlData, 0666)

	if jsonErr != nil {
		spinner.LogMessage("JSON Metadata creation unsuccessful: "+jsonErr.Error(), "fatal")
	}

	if yamlErr != nil {
		spinner.LogMessage("YAML Metadata creation unsuccessful: "+yamlErr.Error(), "fatal")
	}

	// Delete build metadata file grabbed from container to workspace dir
	removeErr := os.Remove(workspaceDir + "/builder_builds.json")
	if removeErr != nil {
		spinner.LogMessage("Couldn't delete temporary log file in workspace: "+removeErr.Error(), "error")
	}

	// Delete docker metadata file created by utils.Metadata()
	removeErr2 := os.Remove(workspaceDir + "/metadata.json")
	if removeErr2 != nil {
		spinner.LogMessage("Couldn't delete temporary log file in workspace: "+removeErr2.Error(), "error")
	}
}

func SaveBuildLogs(buildLogsJSON []byte, path string) {
	dockerLogsJSON, err := os.ReadFile(path + "/docker_logs.json")

	f, err := os.OpenFile(path+"/logs.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		spinner.LogMessage("Cannot create logs.json file: "+err.Error(), "fatal")
	}

	if _, err = f.WriteString(string(buildLogsJSON)); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString("{\"level\":\"info\",\"timestamp\":\"\",\"caller\":\"cmd/docker.go:580\",\"msg\":\"======= Build End =======\"}\n"); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString("{\"level\":\"info\",\"timestamp\":\"\",\"caller\":\"cmd/docker.go:583\",\"msg\":\"====== Docker Start ======\"}\n"); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString(string(buildLogsJSON)); err != nil {
		spinner.LogMessage("unsuccessful write to logs.json file.", "fatal")
	}
	if _, err = f.WriteString(string(dockerLogsJSON)); err != nil {
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

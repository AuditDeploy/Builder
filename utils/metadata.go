package utils

import (
	"Builder/spinner"
	"crypto/sha256"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"encoding/json"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type Artifacts struct {
	name     string
	checksum string
}

// EnvData holds the struct of envs argument
type EnvData struct {
	Key   string
	Value string
}

// AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	ProjectName         string
	ProjectType         string
	ArtifactName        string
	ArtifactChecksums   string
	ArtifactLocation    string
	LogsLocation        string
	UserName            string
	HomeDir             string
	IP                  string
	StartTime           string
	EndTime             string
	GitURL              string
	MasterGitHash       string
	BranchName          string
	AppIcon             string
	ContainerPort       int
	ServicePort         int
	DependsOnCandidates []string
	Envs                []EnvData
}

func Metadata(path string) {
	//Metedata
	projectName := GetName()

	caser := cases.Title(language.English)
	projectType := caser.String(os.Getenv("BUILDER_PROJECT_TYPE"))

	// If we are running the builder docker command we will not have
	// artifact to display in docker metadata so leave as empty
	var artifactName, artifactChecksums string
	if os.Getenv("BUILDER_DOCKER_COMMAND") == "true" {
		artifactName = ""
		artifactChecksums = ""
	} else {
		artifactName = os.Getenv("BUILDER_ARTIFACT_NAMES")
		artifactChecksums = GetArtifactChecksum()
	}

	builderPath, _ := os.Getwd()
	artifactPath := os.Getenv("BUILDER_ARTIFACT_DIR")
	var artifactLocation string
	if os.Getenv("BUILDER_OUTPUT_PATH") != "" {
		artifactLocation = os.Getenv("BUILDER_OUTPUT_PATH")
	} else {
		if os.Getenv("BUILDER_DIR_PATH") != "" {
			if artifactPath[0:1] == "." {
				artifactPath = artifactPath[1:]
				artifactLocation = builderPath + artifactPath
			} else {
				artifactLocation = artifactPath
			}
		} else {
			artifactPath = artifactPath[1:]
			artifactLocation = builderPath + artifactPath
		}
	}

	logsPath := os.Getenv("BUILDER_LOGS_DIR")
	var logsLocation string
	if strings.HasPrefix(logsPath, "./") {
		logsLocation = builderPath + "/" + logsPath[2:] + "/logs.json"
	} else {
		logsLocation = logsPath + "/logs.json"
	}

	ip := GetIPAdress().String()

	userName := GetUserData().Username

	// If on Windows, remove computer ID prefix from returned username
	if runtime.GOOS == "windows" {
		splitString := strings.Split(userName, "\\")
		userName = splitString[1]
	}

	homeDir := GetUserData().HomeDir
	startTime := os.Getenv("BUILD_START_TIME")
	endTime := os.Getenv("BUILD_END_TIME")

	var gitURL = GetRepoURL()
	_, masterGitHash := GitMasterNameAndHash()

	var branchName string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		if os.Getenv("REPO_BRANCH_NAME") != "" {
			branchName = os.Getenv("REPO_BRANCH_NAME")
		} else {
			out, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
			if err != nil {
				spinner.LogMessage("Can't get current branch name.  Please provide it in the builder.yaml.", "info")
				branchName = ""
			} else {
				// remove \n at end of returned branch name before returning
				branchName = strings.TrimSuffix(string(out), "\n")
			}
		}
	} else {
		branchName = os.Getenv("REPO_BRANCH_NAME")
	}

	appIcon := os.Getenv("BUILD_APP_ICON")

	containerPort, _ := strconv.Atoi(os.Getenv("RELEASE_CONTAINER_PORT"))
	servicePort, _ := strconv.Atoi(os.Getenv("RELEASE_SERVICE_PORT"))

	var dependsOnCandidates []string
	if os.Getenv("RELEASE_DEPENDENCIES") == "" {
		dependsOnCandidates = nil
	} else {
		dependsOnCandidates = strings.Split(os.Getenv("RELEASE_DEPENDENCIES"), ",")
	}

	var releaseEnvs []EnvData
	var pairData EnvData
	if os.Getenv("RELEASE_ENVS") != "" {
		envPairs := strings.Split(os.Getenv("RELEASE_ENVS"), ";")
		for _, pair := range envPairs {
			pairArray := strings.Split(pair, ",")
			pairData.Key = pairArray[0]
			pairData.Value = pairArray[1]
			releaseEnvs = append(releaseEnvs, pairData)
		}
	} else {
		releaseEnvs = nil
	}

	//Contains a collection of files with user's metadata
	userMetaData := AllMetaData{
		ProjectName:         projectName,
		ProjectType:         projectType,
		ArtifactName:        artifactName,
		ArtifactChecksums:   artifactChecksums,
		ArtifactLocation:    artifactLocation,
		LogsLocation:        logsLocation,
		UserName:            userName,
		HomeDir:             homeDir,
		IP:                  ip,
		StartTime:           startTime,
		EndTime:             endTime,
		GitURL:              gitURL,
		MasterGitHash:       masterGitHash,
		BranchName:          branchName,
		AppIcon:             appIcon,
		ContainerPort:       containerPort,
		ServicePort:         servicePort,
		DependsOnCandidates: dependsOnCandidates,
		Envs:                releaseEnvs,
	}

	OutputMetadata(path, &userMetaData)

}

// GetUserData return username and userdir
func GetUserData() *user.User {
	user, err := user.Current()
	if err != nil {
		panic(err)

	}

	return user

}

// GetIPAdress Get preferred outbound ip of this machine
func GetIPAdress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		spinner.LogMessage("could not connect to outbound ip: "+err.Error(), "fatal")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// OutputJSONall  outputs allMetaData struct in JSON format
func OutputMetadata(path string, allData *AllMetaData) {
	yamlData, _ := yaml.Marshal(allData)
	jsonData, _ := json.Marshal(allData)

	err := os.WriteFile(path+"/metadata.json", jsonData, 0666)
	err2 := os.WriteFile(path+"/metadata.yaml", yamlData, 0666)

	if err != nil {
		spinner.LogMessage("JSON Metadata creation unsuccessful.", "fatal")
	}

	if err2 != nil {
		spinner.LogMessage("YAML Metadata creation unsuccessful.", "fatal")
	}
}

// Gets the name of the repo's master branch and its hash
func GitMasterNameAndHash() (string, string) {
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	currentDir, _ := os.Getwd()
	var dirToRunIn string
	if os.Getenv("BUILDER_COMMAND") == "true" || os.Getenv("BUILDER_DOCKER_COMMAND") == "true" {
		dirToRunIn = currentDir
	} else {
		dirToRunIn, _ = filepath.Abs(hiddenDir)
	}

	//outputs the name of the master branch
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short")
	if os.Getenv("BUILDER_COMMAND") != "true" {
		cmd.Dir = dirToRunIn
	}
	output, err := cmd.Output()
	if err != nil {
		// Can't find master branch name so return undefined
		return "undefined", "undefined"
	}
	formattedOutput := strings.TrimSuffix(string(output), "\n")
	masterBranchName := formattedOutput[strings.LastIndex(formattedOutput, "/")+1:]

	//outputs the hash of the provided branch
	cmd = exec.Command("git", "rev-parse", masterBranchName)
	if os.Getenv("BUILDER_COMMAND") != "true" {
		cmd.Dir = dirToRunIn
	}
	hashOutput, hashErr := cmd.Output()
	if hashErr != nil {
		return "undefined", "undefined"
	}
	masterBranchHash := strings.TrimSuffix(string(hashOutput), "\n")

	return masterBranchName, masterBranchHash
}

func GetArtifactChecksum() string {
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	files, err := os.ReadDir(artifactDir)
	if err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	var checksum, checksums string
	var checksumsArray []Artifacts
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, file := range files {
			if file.Name() != "metadata.json" && file.Name() != "metadata.yaml" {
				// Get checksum of artifact
				artifact, err := os.ReadFile(artifactDir + "/" + file.Name())
				if err != nil {
					spinner.LogMessage(err.Error(), "fatal")
				}

				sum := sha256.Sum256(artifact)
				checksum = fmt.Sprintf("%x", sum)

				var artifactObj Artifacts
				artifactObj.name = file.Name()
				artifactObj.checksum = checksum

				checksumsArray = append(checksumsArray, artifactObj)
				checksums = fmt.Sprintf("%+v", checksumsArray)
			}
		}
	}()
	wg.Wait()

	return checksums
}

func GetBuildID() string {
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	var checksum string

	// Get checksum of metadata.json
	metadata, err := os.ReadFile(artifactDir + "/metadata.json")
	if err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	sum := sha256.Sum256(metadata)
	checksum = fmt.Sprintf("%x", sum)

	// Only return first 10 char of sum
	return checksum[0:9]
}

func StoreBuildMetadataLocally() {
	// Read in build JSON data from build artifact directory
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	metadataJSON, err := os.ReadFile(artifactDir + "/metadata.json")
	if err != nil {
		spinner.LogMessage("Cannot find metadata.json file: "+err.Error(), "fatal")
	}

	// Unmarshal json data so we can add buildID
	var metadataFormat map[string]interface{}
	json.Unmarshal(metadataJSON, &metadataFormat)
	metadataFormat["BuildID"] = GetBuildID()

	updatedMetadataJSON, err := json.Marshal(metadataFormat)
	if err != nil {
		spinner.LogMessage("Cannot marshal metadata: "+err.Error(), "fatal")
	}

	// Check if builds.json exists and append to it, if not, create it
	textToAppend := string(updatedMetadataJSON) + ",\n"

	var pathToBuildsJSON string

	if runtime.GOOS == "windows" {
		appDataDir := os.Getenv("LOCALAPPDATA")
		if appDataDir == "" {
			appDataDir = os.Getenv("APPDATA")
		}

		pathToBuildsJSON = appDataDir + "/Builder/builds.json"
	} else {
		user, _ := user.Current()
		homeDir := user.HomeDir

		pathToBuildsJSON = homeDir + "/.builder/builds.json"
	}

	buildsFile, err := os.OpenFile(pathToBuildsJSON, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		spinner.LogMessage("Could not create builds.json file: "+err.Error(), "fatal")
	}
	if _, err := buildsFile.WriteString(textToAppend); err != nil {
		spinner.LogMessage("Could not write to builds.json file: "+err.Error(), "fatal")
	}
	if err := buildsFile.Close(); err != nil {
		spinner.LogMessage("Could not close builds.json file: "+err.Error(), "fatal")
	}
}

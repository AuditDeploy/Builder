package utils

import (
	"Builder/spinner"
	"crypto/sha256"
	"fmt"
	"runtime"

	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

func Metadata(path string) {
	//Metedata
	projectName := GetName()

	caser := cases.Title(language.English)
	projectType := caser.String(os.Getenv("BUILDER_PROJECT_TYPE"))

	artifactName := os.Getenv("BUILDER_ARTIFACT_NAMES")
	artifactChecksums := GetArtifactChecksum()
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
	var masterGitHash string
	if os.Getenv("BUILDER_COMMAND") != "true" {
		_, masterGitHash = GitHashAndName()
	}

	branchName := os.Getenv("REPO_BRANCH_NAME")

	//Contains a collection of files with user's metadata
	userMetaData := AllMetaData{
		ProjectName:       projectName,
		ProjectType:       projectType,
		ArtifactName:      artifactName,
		ArtifactChecksums: artifactChecksums,
		ArtifactLocation:  artifactLocation,
		LogsLocation:      logsLocation,
		UserName:          userName,
		HomeDir:           homeDir,
		IP:                ip,
		StartTime:         startTime,
		EndTime:           endTime,
		GitURL:            gitURL,
		MasterGitHash:     masterGitHash,
		BranchName:        branchName}

	OutputMetadata(path, &userMetaData)
}

// AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	ProjectName       string
	ProjectType       string
	ArtifactName      string
	ArtifactChecksums string
	ArtifactLocation  string
	LogsLocation      string
	UserName          string
	HomeDir           string
	IP                string
	StartTime         string
	EndTime           string
	GitURL            string
	MasterGitHash     string
	BranchName        string
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

	err := ioutil.WriteFile(path+"/metadata.json", jsonData, 0666)
	err2 := ioutil.WriteFile(path+"/metadata.yaml", yamlData, 0666)

	if err != nil {
		spinner.LogMessage("JSON Metadata creation unsuccessful.", "fatal")
	}

	if err2 != nil {
		spinner.LogMessage("YAML Metadata creation unsuccessful.", "fatal")
	}
}

// GitHas gets the latest git commit id in a repo
func GitHashAndName() ([]string, string) {
	//Get repoURL
	repo := GetRepoURL()

	//outputs all the commits of the clone repo
	output, _ := exec.Command("git", "ls-remote", repo).Output()

	//stringify output - []byte to string
	stringGitHashAndName := string(output)
	// fmt.Println(stringGitHash)

	//return an array with all the git commit hashs
	arrayGitHashAndName := strings.Split(stringGitHashAndName, "\n")

	//gets the hash of type []string of master branch
	masterHashStringArray := strings.Fields(arrayGitHashAndName[0])
	masterHash := masterHashStringArray[0]

	return arrayGitHashAndName, masterHash[0:7]
}

type Artifacts struct {
	name     string
	checksum string
}

func GetArtifactChecksum() string {
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	files, err := os.ReadDir(artifactDir)
	if err != nil {
		spinner.LogMessage(err.Error(), "fatal")
	}

	var checksumsArray []Artifacts
	var checksum string
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
		}
	}
	checksums := fmt.Sprintf("%+v", checksumsArray)

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
	if _, err := buildsFile.Write([]byte(textToAppend)); err != nil {
		spinner.LogMessage("Could not write to builds.json file: "+err.Error(), "fatal")
	}
	if err := buildsFile.Close(); err != nil {
		spinner.LogMessage("Could not close builds.json file: "+err.Error(), "fatal")
	}
}

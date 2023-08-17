package utils

import (
	"bytes"
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

	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var BuilderLog = zap.S()

func Metadata(path string) {
	//Metedata
	projectName := GetName()

	caser := cases.Title(language.English)
	projectType := caser.String(os.Getenv("BUILDER_PROJECT_TYPE"))

	artifactName := os.Getenv("BUILDER_ARTIFACT_STAMP")
	artifactChecksum := GetArtifactChecksum()
	var artifactLocation string
	if os.Getenv("BUILDER_OUTPUT_PATH") != "" {
		artifactLocation = os.Getenv("BUILDER_OUTPUT_PATH")
	} else {
		builderPath, _ := os.Getwd()
		artifactRelativePath := os.Getenv("BUILDER_ARTIFACT_DIR")
		artifactDir := artifactRelativePath[1:]

		artifactLocation = builderPath + artifactDir
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

	var masterGitHash string
	if os.Getenv("BUILDER_COMMAND") != "true" {
		_, masterGitHash = GitHashAndName()
	}

	branchName := os.Getenv("REPO_BRANCH_NAME")

	//Contains a collection of files with user's metadata
	userMetaData := AllMetaData{
		ProjectName:      projectName,
		ProjectType:      projectType,
		ArtifactName:     artifactName,
		ArtifactChecksum: artifactChecksum,
		ArtifactLocation: artifactLocation,
		UserName:         userName,
		HomeDir:          homeDir,
		IP:               ip,
		StartTime:        startTime,
		EndTime:          endTime,
		MasterGitHash:    masterGitHash,
		BranchName:       branchName}

	OutputMetadata(path, &userMetaData)
}

// AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	ProjectName      string
	ProjectType      string
	ArtifactName     string
	ArtifactChecksum string
	ArtifactLocation string
	UserName         string
	HomeDir          string
	IP               string
	StartTime        string
	EndTime          string
	MasterGitHash    string
	BranchName       string
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
		BuilderLog.Fatalf("could not connect to outbound ip", err)
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
		BuilderLog.Fatal("JSON Metadata creation unsuccessful.")
	}

	if err2 != nil {
		BuilderLog.Fatal("YAML Metadata creation unsuccessful.")
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

func GetArtifactChecksum() string {
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	files, err := os.ReadDir(artifactDir)
	if err != nil {
		BuilderLog.Fatal(err)
	}

	var checksum string
	for _, file := range files {
		if file.Name() != "metadata.json" && file.Name() != "metadata.yaml" {
			// Get checksum of artifact
			artifact, err := os.ReadFile(artifactDir + "/" + file.Name())
			if err != nil {
				BuilderLog.Fatal(err)
			}

			sum := sha256.Sum256(artifact)
			checksum = fmt.Sprintf("%x", sum)
		}
	}

	return checksum
}

func GetBuildID() string {
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")
	var checksum string

	// Get checksum of metadata.json
	metadata, err := os.ReadFile(artifactDir + "/metadata.json")
	if err != nil {
		BuilderLog.Fatal(err)
	}

	sum := sha256.Sum256(metadata)
	checksum = fmt.Sprintf("%x", sum)

	// Only return first 10 char of sum
	return checksum[0:9]
}

// AllMetaData holds the stuct of all the arguments
type MetadataFormat struct {
	ProjectName      string `json:"ProjectName"`
	ProjectType      string `json:"ProjectType"`
	ArtifactName     string `json:"ArtifactName"`
	ArtifactChecksum string `json:"ArtifactChecksum"`
	ArtifactLocation string `json:"ArtifactLocation"`
	UserName         string `json:"UserName"`
	HomeDir          string `json:"HomeDir"`
	IP               string `json:"IP"`
	StartTime        string `json:"StartTime"`
	EndTime          string `json:"EndTime"`
	MasterGitHash    string `json:"MasterGitHash"`
	BranchName       string `json:"BranchName"`
}

func StoreBuildMetadataLocally() {
	// Read in build JSON data from build artifact directory
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	metadataJSON, err := os.ReadFile(artifactDir + "/metadata.json")
	if err != nil {
		var _, errb bytes.Buffer
		BuilderLog.Fatalf("Cannot find metadata.json file", errb)
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
		BuilderLog.Fatalf("Could not create builds.json file", err)
	}
	if _, err := buildsFile.Write([]byte(textToAppend)); err != nil {
		BuilderLog.Fatalf("Could not write to builds.json file", err)
	}
	if err := buildsFile.Close(); err != nil {
		BuilderLog.Fatalf("Could not close builds.json file", err)
	}
}

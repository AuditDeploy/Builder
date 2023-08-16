package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var BuilderLog = zap.S()

func Metadata(path string) {
	//Metedata
	timestamp := time.Now().Format(time.RFC850)
	ip := GetIPAdress().String()
	userName := GetUserData().Username
	homeDir := GetUserData().HomeDir

	var masterGitHash, branchHash, branchName string
	if os.Getenv("BUILDER_COMMAND") != "true" {
		_, masterGitHash, branchHash, branchName = GitHashAndName()
	}

	//Contains a collection of fileds with user's metadata
	userMetaData := AllMetaData{
		UserName:      userName,
		HomeDir:       homeDir,
		IP:            ip,
		Timestamp:     timestamp,
		MasterGitHash: masterGitHash,
		BranchName:    branchName,
		BranchHash:    branchHash}

	OutputMetadata(path, &userMetaData)
}

// AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	UserName      string
	HomeDir       string
	IP            string
	Timestamp     string
	MasterGitHash string
	BranchName    string
	BranchHash    string
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
func GitHashAndName() ([]string, string, string, string) {
	//Get repoURL
	repo := GetRepoURL()

	//outputs all the commits of the clone repo
	output, _ := exec.Command("git", "ls-remote", repo).Output()

	//stringify output - []byte to string
	stringGitHashAndName := string(output)
	// fmt.Println(stringGitHash)

	//return an array with all the git commit hashs
	arrayGitHashAndName := strings.Split(stringGitHashAndName, "\n")
	branchExists, branchNameAndHash := BranchNameExists(arrayGitHashAndName)

	//gets the hash of type []string of master branch
	masterHashStringArray := strings.Fields(arrayGitHashAndName[0])
	masterHash := masterHashStringArray[0]

	if branchExists {
		//gets hash and name of type []string of a specific branch
		branchHash := strings.Fields(branchNameAndHash)[0]
		branchName := strings.Fields(branchNameAndHash)[1]
		return arrayGitHashAndName, masterHash[0:7], branchHash[0:7], branchName
	} else {
		return arrayGitHashAndName, masterHash[0:7], "", ""
	}

}

func BranchNameExists(branches []string) (bool, string) {
	branchExists := false
	var branchNameAndHash string

	_, clonedBranchName := CloneBranch()

	for _, branch := range branches {
		nameSlice := strings.Split(branch, "/")
		sliceLen := len(nameSlice)
		if branch[strings.LastIndex(branch, "/")+1:] == clonedBranchName || (sliceLen > 2 && (nameSlice[sliceLen-2]+"/"+nameSlice[sliceLen-1] == clonedBranchName)) {
			branchExists = true
			branchNameAndHash = branch
		}
	}
	return branchExists, branchNameAndHash
}

func StoreBuildMetadataLocally() {
	// Read in build JSON data from build artifact directory
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	metadataJSON, err := os.ReadFile(artifactDir + "/metadata.json")
	if err != nil {
		var _, errb bytes.Buffer
		BuilderLog.Fatalf("Cannot find metadata.json file", errb)
	}

	// Check if builds.json exists and append to it, if not, create it
	textToAppend := string(metadataJSON) + ",\n"

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

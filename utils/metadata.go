package utils

import (
	"Builder/utils/log"

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

	artifactName := os.Getenv("BUILDER_ARTIFACT_STAMP")
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
	homeDir := GetUserData().HomeDir
	startTime := os.Getenv("BUILD_START_TIME")
	endTime := os.Getenv("BUILD_END_TIME")

	var masterGitHash string
	if os.Getenv("BUILDER_COMMAND") != "true" {
		_, masterGitHash = GitHashAndName()
	}

	branchName := os.Getenv("REPO_BRANCH_NAME")
	buildHash := os.Getenv("BUILD_HASH")

	//Contains a collection of files with user's metadata
	userMetaData := AllMetaData{
		ProjectName:      projectName,
		ProjectType:      projectType,
		ArtifactName:     artifactName,
		ArtifactLocation: artifactLocation,
		UserName:         userName,
		HomeDir:          homeDir,
		IP:               ip,
		StartTime:        startTime,
		EndTime:          endTime,
		MasterGitHash:    masterGitHash,
		BranchName:       branchName,
		BuildHash:        buildHash}

	OutputMetadata(path, &userMetaData)
}

// AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	ProjectName      string
	ProjectType      string
	ArtifactName     string
	ArtifactLocation string
	UserName         string
	HomeDir          string
	IP               string
	StartTime        string
	EndTime          string
	MasterGitHash    string
	BranchName       string
	BuildHash        string
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
		log.Fatal("could not connect to outbound ip", err)
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
		log.Fatal("JSON Metadata creation unsuccessful.")
	}

	if err2 != nil {
		log.Fatal("YAML Metadata creation unsuccessful.")
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

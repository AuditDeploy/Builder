package utils

import (
	"Builder/logger"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func Metadata() {
	//Metedata
	timestamp := time.Now().Format(time.RFC850)
	ip := GetIPAdress().String()
	userName := GetUserData().Username
	homeDir := GetUserData().HomeDir
	_, masterGitHash, branchHash, branchName := GitHashAndName()

	//Contains a collection of fileds with user's metadata
	userMetaData := AllMetaData{
		UserName:      userName,
		HomeDir:       homeDir,
		IP:            ip,
		Timestamp:     timestamp,
		MasterGitHash: masterGitHash,
		BranchName:    branchName,
		BranchHash:    branchHash}

	OutputMetadata(&userMetaData)

}

//AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	UserName      string
	HomeDir       string
	IP            string
	Timestamp     string
	MasterGitHash string
	BranchName    string
	BranchHash    string
}

//GetUserData return username and userdir
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
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//OutputJSONall  outputs allMetaData struct in JSON format
func OutputMetadata(allData *AllMetaData) {
	parentDir := os.Getenv("BUILDER_PARENT_DIR")

	yamlData, _ := yaml.Marshal(allData)
	jsonData, _ := json.Marshal(allData)

	err := ioutil.WriteFile(parentDir+"/metedata.json", jsonData, 0644)
	err2 := ioutil.WriteFile(parentDir+"/metedata.yml", yamlData, 0644)

	if err != nil {
		logger.ErrorLogger.Println("JSON Metadata creation unsuccessful.")
		panic(err)
	}

	if err2 != nil {
		logger.ErrorLogger.Println("YAML Metadata creation unsuccessful.")
		panic(err2)
	}
}

//GitHas gets the latest git commit id in a repo
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
		if branch[strings.LastIndex(branch, "/")+1:] == clonedBranchName {
			branchExists = true
			branchNameAndHash = branch
		}
	}

	return branchExists, branchNameAndHash
}

package utils

import (
	"Builder/logger"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"time"

	"gopkg.in/yaml.v2"
)

func Metadata(path string) {
	//Metedata
	timestamp := time.Now().Format(time.RFC850)
	ip := GetIPAdress().String()
	userName := GetUserData().Username
	homeDir := GetUserData().HomeDir

	var masterGitHash, branchGitHash, branchName string
	if os.Getenv("BUILDER_COMMAND") != "true" {
		masterGitHash = GetMasterGitHash()
		branchGitHash = GetBranchGitHash()
		branchName = GetBranchName()
	}

	//Contains a collection of fileds with user's metadata
	userMetaData := AllMetaData{
		UserName:      userName,
		HomeDir:       homeDir,
		IP:            ip,
		Timestamp:     timestamp,
		MasterGitHash: masterGitHash,
		BranchName:    branchName,
		BranchGitHash: branchGitHash,
	}

	OutputMetadata(path, &userMetaData)

}

//AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	UserName      string
	HomeDir       string
	IP            string
	Timestamp     string
	MasterGitHash string
	BranchName    string
	BranchGitHash string
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
func OutputMetadata(path string, allData *AllMetaData) {
	yamlData, _ := yaml.Marshal(allData)
	jsonData, _ := json.Marshal(allData)

	err := ioutil.WriteFile(path+"/metadata.json", jsonData, 0666)
	err2 := ioutil.WriteFile(path+"/metadata.yaml", yamlData, 0666)

	if err != nil {
		logger.ErrorLogger.Println("JSON Metadata creation unsuccessful.")
		panic(err)
	}

	if err2 != nil {
		logger.ErrorLogger.Println("YAML Metadata creation unsuccessful.")
		panic(err2)
	}
}

package utils

import (
	"encoding/json"
	"fmt"
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
	gitHash := GitHash()

	//Contains a collection of fileds with user's metadata
	userMetaData := AllMetaData{
		UserName:  userName,
		HomeDir:   homeDir,
		IP:        ip,
		Timestamp: timestamp,
		GitHash:   gitHash}

	OutputMetadata(&userMetaData)

}

//AllMetaData holds the stuct of all the arguments
type AllMetaData struct {
	UserName  string
	HomeDir   string
	IP        string
	Timestamp string
	GitHash   string
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
		panic(err)
	}

	if err2 != nil {
		panic(err2)
	}
}

//GitHas gets the latest git commit id in a repo
func GitHash() string {
	//Get repo
	args := os.Args[1:]

	//grab URL first
	var repo string
	for i, v := range args {
		if v == "--repo" || v == "-r" {
			if len(args) <= i+1 {
				fmt.Println("No Repo Url Provided")
				os.Exit(1)

			} else {
				repo = args[i+1]
			}
		}
	}
	if repo == "" {
		fmt.Println("No Repo Url Provided")
		os.Exit(1)
	}

	//outputs all the commits of the clone repo
	output, _ := exec.Command("git", "ls-remote", repo, "refs/heads/master").Output()

	//stringify output - []byte to string
	stringGitHash := string(output)

	//return an array with all the git commit hashs
	arrayGitHashs := strings.Split(stringGitHash, "\n")

	//gets the hash of type []string
	hashStringArray := strings.Fields(arrayGitHashs[0])

	return hashStringArray[0]

}

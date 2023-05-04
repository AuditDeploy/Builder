package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
)

func PushRepo() {

	var Cfg BuilderConfig
	var envCfg EnvConfig

	env := envCfg.ReadConfigFile("app.env", "env")

	auth := &http.BasicAuth{
		Username: env.Username,
		Password: env.Password,
	}

	repo := GetRepoURL()
	path := filepath.Join("tmp", "foo")

	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})

	defer os.RemoveAll(path)

	if err != nil {
		log.Fatal(err)
	}

	w, err := r.Worktree()

	if err != nil {
		log.Fatal(err)
	}

	configType := "json"
	name := GetRepoName(repo)
	currentDirectory, _ := os.Getwd()
	localPath := filepath.Join(currentDirectory, "configs")
	tmpPath := filepath.Join("tmp", "foo")
	filename, v := CreateConfigFile(name+"_builder", configType)
	builderFile := filepath.Join(tmpPath, name+"_builder"+configType)
	var allSettings map[string]interface{}

	// if v == nil; file already exists; get config from configuration slice created in main
	if v == nil {
		_, err := os.Stat(builderFile)

		// builder json file dosen't exist in cloned repo use local builder json file
		if os.IsNotExist(err) {
			cfgs := RetrieveConfis()

			for _, c := range cfgs {
				if c.ConfigFileUsed() == filename {
					v = c
				}
			}
			allSettings = Cfg.ReadConfigFile(name+"_builder.json", configType, localPath, v)
		} else {
			// use builder json file from repo
			v = viper.New()
			allSettings = Cfg.ReadConfigFile(name+"_builder.json", configType, tmpPath, v)
		}
	} else {
		// if v != nil; local builder json file dosen't exist
		_, err := os.Stat(builderFile)

		// builder json file exists in cloned repo read from repo builder json file
		if !(os.IsNotExist(err)) {
			allSettings = Cfg.ReadConfigFile(name+"_builder.json", configType, tmpPath, v)
		} else {
			// builder json file dosen't exist in repo read from local builder json file just created
			allSettings = Cfg.ReadConfigFile(name+"_builder.json", configType, localPath, v)
		}
	}

	dt := time.Now()
	dtformatted := dt.Format("2006-01-02 15:04:05")

	allSettings["date"] = dtformatted

	val, ok := allSettings["buildnumber"]

	if !ok {
		allSettings["buildnumber"] = 1
	} else {
		buildNumber := val.(int64)
		allSettings["buildnumber"] = buildNumber + 1
	}

	WriteConfigFile(name+"_builder."+configType, configType, localPath, &allSettings, v)
	exec.Command("cp", filename, filepath.Join("tmp", "foo")).Run()

	_, err = w.Add(name + "_builder." + configType)
	if err != nil {
		fmt.Println("Add error: ", err)
	}

	val = allSettings["buildnumber"]
	buildNumber := val.(int64)

	w.Commit("Added new build number "+strconv.FormatInt(int64(buildNumber), 10), &git.CommitOptions{})

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		log.Fatal("Error pushing to repo: ", err)
	}

	fmt.Println("Remote updated ", filename)
}

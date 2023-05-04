package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func CreateConfigFile(name string, configType string) (string, *viper.Viper) {
	currentDir, _ := os.Getwd()
	configDir := filepath.Join(currentDir, "configs")

	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create configs directory: ", err)
		}
	}

	// check if file exists
	filename := filepath.Join(configDir, name+"."+configType)
	var _, err = os.Stat(filename)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal("Error creating config file: ", err)
		}
		defer file.Close()
		fmt.Println("File Created Successfully", filename)

		v := viper.New()

		addConfig(v)

		return filename, v
	}

	return filename, nil
}

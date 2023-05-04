package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func WriteConfigFile(filename string, configType string, configPath string, settings *map[string]interface{}, v *viper.Viper) {
	file := strings.Split(filename, ".")[0]

	v.SetConfigName(file)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	fmt.Println("settings: ", *settings)

	for key, value := range *settings {
		v.Set(key, value)
	}

	err := v.WriteConfigAs(filepath.Join(configPath, filename))

	if err != nil {
		log.Fatal("Error writing config file ", err)
	}
}

package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Autopush string `mapstructure:"autopush"`
}

type BuilderConfig struct {
	Date        string `mapstructure:"date"`
	BuildNumber int64  `mapstructure:"buildnumber"`
}

type EnvConfig struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

var viperConfigs = make([]*viper.Viper, 0)

func (cfg Config) ReadConfigFile(filename string, configType string, configPath string, v *viper.Viper) map[string]interface{} {
	file := strings.Split(filename, ".")[0]

	v.SetConfigName(file)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	err := v.ReadInConfig()
	if err != nil && !strings.Contains(err.Error(), "unexpected end of") {
		log.Fatal("Error reading config file: "+filename+" ", err)
	}

	err = v.Unmarshal(&cfg)

	if err != nil {
		fmt.Println("Unmarshalling failed!")
	}

	m := make(map[string]interface{})

	values := reflect.ValueOf(cfg)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		key := strings.ToLower(types.Field(i).Name)
		value := values.Field(i).Interface()
		if !(value == "" || value == 0 || value == 0.0 || value == nil || value == false) {
			m[key] = values.Field(i).Interface()
		}
	}

	return m
}

func (builderCfg BuilderConfig) ReadConfigFile(filename string, configType string, configPath string, v *viper.Viper) map[string]interface{} {
	file := strings.Split(filename, ".")[0]

	v.SetConfigName(file)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	err := v.ReadInConfig()
	if err != nil && !strings.Contains(err.Error(), "unexpected end of") {
		log.Fatal("Error reading config file: "+filename+" ", err)
	}

	err = v.Unmarshal(&builderCfg)

	if err != nil {
		fmt.Println("Unmarshalling failed!")
	}

	m := make(map[string]interface{})

	values := reflect.ValueOf(builderCfg)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		key := strings.ToLower(types.Field(i).Name)
		value := values.Field(i).Interface()
		if !(value == "" || value == 0 || value == 0.0 || value == nil || value == false) {
			m[key] = values.Field(i).Interface()
		}
	}

	return m
}

func (envCfg EnvConfig) ReadConfigFile(filename string, configType string) *EnvConfig {
	file := strings.Split(filename, ".")[0]

	viper.SetConfigName(file)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil && !strings.Contains(err.Error(), "unexpected end of") {
		log.Fatal("Error reading config file: "+filename+" ", err)
	}

	err = viper.Unmarshal(&envCfg)

	if err != nil {
		log.Fatal("Unmarshalling failed! ", err)
	}

	return &envCfg
}

func RetrieveConfis() []*viper.Viper {
	return viperConfigs
}

func addConfig(cfg *viper.Viper) {
	viperConfigs = append(viperConfigs, cfg)
}

func InitConfig() {
	currentDir, _ := os.Getwd()
	configDir := filepath.Join(currentDir, "configs")

	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		// configDir dosen't exist return from function
		return
	}

	files, err := os.ReadDir(configDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileSplit := strings.Split(file.Name(), ".")
		fileName := fileSplit[0]
		fileExt := fileSplit[1]

		v := viper.New()
		v.SetConfigName(fileName)
		v.SetConfigType(fileExt)
		v.AddConfigPath(configDir)
		v.ReadInConfig()

		addConfig(v)
	}
}

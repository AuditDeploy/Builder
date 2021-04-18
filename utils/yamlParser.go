package utils

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func YamlParser() {
	// Initial declaration
 m := map[string]interface{}{
	"key": "value",
}
// Dynamically add a sub-map
 m["sub"] = map[string]interface{}{
	"deepKey": "deepValue",
}
// returns map
 var f interface{}

//takes yaml path and read file
	source, err := ioutil.ReadFile("./tempRepo/builder.yaml")
	if err != nil {
		removeTempDir()
		log.Fatal(err)
	}

	//unpacks yaml file in a map int{}
	err = yaml.Unmarshal([]byte(source), &f) 
	if err != nil {
		log.Printf("error: %v", err)
	}

	//pass map int{} to callback that sets env vars
	ConfigEnvs(f)

	removeTempDir()
}

func removeTempDir() {
	//delete tempRepo dir
	err := os.RemoveAll("./tempRepo")
	if err != nil {
			log.Fatal(err)
	}
}
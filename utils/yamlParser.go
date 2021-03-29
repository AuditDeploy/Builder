package utils

import (
	"fmt"
	"io/ioutil"
	"log"

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

 //~~~~FOR FUTURE, TO MAKE DYNAMIC~~~~
//  workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
//  yamlPath, _ := exec.Command("find", workspaceDir, "-name", "builder.yaml").CombinedOutput()
//  yamlString := string(yamlPath)

//takes yaml path and read file
	source, err := ioutil.ReadFile("builder.yaml")
	if err != nil {
		log.Fatal(err)
	}

	//print raw yaml
	fmt.Printf("File contents: %s", source)

	//unpacks yaml file in a map 
	err = yaml.Unmarshal([]byte(source), &f) 
	if err != nil {
		log.Printf("error: %v", err)
	}

	//print map of yaml
	fmt.Printf("%+v\n", f)
}
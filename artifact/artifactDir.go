package artifact

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func ArtifactDir() {
	var dirPath string
	// if os.Getenv("builder_COMMAND") == "true" {
	// 	path, _ := os.Getwd()
	// 	if strings.Contains(path, "workspace") {
	// 		dirPath = strings.TrimRight(path, "\\workspace")
	// 	} else if strings.Contains(path, "workspace") && strings.Contains(path, "temp") {
	// 		dirPath = strings.TrimRight(path, "\\temp")
	// 	}
	// } else {
	dirPath = os.Getenv("builder_PARENT_DIR")
	// }

	currentTime := time.Now().Unix()
	artifactStamp := "artifact_" + strconv.FormatInt(currentTime, 10)
	os.Setenv("builder_ARTIFACT_STAMP", artifactStamp)
	artifactDir := dirPath + "/" + artifactStamp

	err := os.Mkdir(artifactDir, 0755)
	//should return nil once directory is made, if not, throw err
	if err != nil {
		log.Fatal(err)
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("builder_ARTIFACT_DIR")
	if !present {
		os.Setenv("builder_ARTIFACT_DIR", artifactDir)
	} else {
		fmt.Println("builder_ARTIFACT_DIR", val)
	}
}

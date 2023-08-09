package artifact

import (
	"Builder/utils/log"
	"fmt"
	"os"
	"strconv"
	"time"
)

func ArtifactDir(startTime string) {
	var dirPath string
	// if os.Getenv("BUILDER_COMMAND") == "true" {
	// 	path, _ := os.Getwd()
	// 	if strings.Contains(path, "workspace") {
	// 		dirPath = strings.TrimRight(path, "\\workspace")
	// 	} else if strings.Contains(path, "workspace") && strings.Contains(path, "temp") {
	// 		dirPath = strings.TrimRight(path, "\\temp")
	// 	}
	// } else {

	dirPath = os.Getenv("BUILDER_PARENT_DIR")
	dirName := os.Getenv("BUILDER_DIR_NAME")
	// }

	parsedStartTime, _ := time.Parse(time.RFC850, startTime)
	timeBuildStarted := parsedStartTime.Unix()
	artifactStamp := dirName + "_artifact_" + strconv.FormatInt(timeBuildStarted, 10)
	os.Setenv("BUILDER_ARTIFACT_STAMP", artifactStamp)
	artifactDir := dirPath + "/" + artifactStamp

	err := os.Mkdir(artifactDir, 0755)
	//should return nil once directory is made, if not, throw err
	if err != nil {
		log.Fatal("failed to make artifact directory", err)
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_ARTIFACT_DIR")
	if !present {
		os.Setenv("BUILDER_ARTIFACT_DIR", artifactDir)
	} else {
		fmt.Println("BUILDER_ARTIFACT_DIR", val)
	}
}

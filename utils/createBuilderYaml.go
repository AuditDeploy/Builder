package utils

import (
	"log"
	"os"
)

func CreateBuilderYaml() {
	workSpace := os.Getenv("BUILDER_WORKSPACE_DIR")

	file, err := os.OpenFile(workSpace+"/"+"builder.yaml", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

}

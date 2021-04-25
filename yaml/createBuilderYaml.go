package yaml

import (
	"log"
	"os"
)

func CreateBuilderYaml(filePath string) {

	file, err := os.OpenFile(filePath+"/"+"builder.yaml", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

}


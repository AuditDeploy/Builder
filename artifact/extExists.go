package artifact

import (
	"Builder/utils/log"
	"fmt"
	"os"
	"path/filepath"
)

// find file with extension and return file name
func ExtExistsFunction(dirPath string, ext string) (bool, string) {
	found := false
	d, err := os.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		log.Fatal("could not find dirpath %v", dirPath, err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		log.Fatal("could not read directory", err)
		os.Exit(1)
	}
	var fileName string

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ext {
				fileName = file.Name()
				found = true
			}
		}
	}
	return found, fileName
}

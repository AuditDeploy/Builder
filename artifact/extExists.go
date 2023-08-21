package artifact

import (
	"Builder/spinner"
	"Builder/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// find file with extension and return file name
func ExtExistsFunction(dirPath string, ext string) (bool, string) {
	found := false
	d, err := os.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		spinner.LogMessage("could not find dirpath "+dirPath+": "+err.Error(), "fatal")
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		spinner.LogMessage("could not read directory: "+err.Error(), "fatal")
		os.Exit(1)
	}
	var fileName string

	for _, file := range files {
		if file.Mode().IsRegular() {
			if ext != "executable" {
				if filepath.Ext(file.Name()) == ext {
					fileName = file.Name()
					found = true
				}
			} else {
				if file.Mode()&0111 != 0 && file.Name() == strings.TrimSuffix(utils.GetName(), ".git") {
					fileName = file.Name()
					found = true
				}
			}

		}
	}
	return found, fileName
}

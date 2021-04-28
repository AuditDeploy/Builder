package artifact

import (
	"Builder/logger"
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Npm creates zip from files passed in as arg
func ZipArtifactDir() {
	// parentDir := os.Getenv("BUILDER_PARENT_DIR")
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	artifactZip := artifactDir+".zip"

	// CreateZip temp dir.
	outFile, err := os.Create(artifactZip)
	if err != nil {
		 log.Fatal(err)
	}
	
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from artifact dir to the artifact zip.
	addNpmFiles(w, artifactDir+"/", "")

	err = w.Close()
	if err != nil {
		logger.ErrorLogger.Println("Npm project failed to compile.")
		 log.Fatal(err)
	}
}

//recursively add files
func addNpmFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
			fmt.Println(err)
	}

	for _, file := range files {
			if !file.IsDir() {
					dat, err := ioutil.ReadFile(basePath + file.Name())
					if err != nil {
							fmt.Println(err)
					}

					// Add some files to the archive.
					f, err := w.Create(baseInZip + file.Name())
					if err != nil {
							fmt.Println(err)
					}
					_, err = f.Write(dat)
					if err != nil {
							fmt.Println(err)
					}
			} else if file.IsDir() {
					// Recurse
					newBase := basePath + file.Name() + "/"
					addNpmFiles(w, newBase, baseInZip  + file.Name() + "/")
			}
	}
}
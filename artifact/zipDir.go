package artifact

import (
	"Builder/spinner"
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

// Npm creates zip from files passed in as arg
func ZipArtifactDir() {
	// parentDir := os.Getenv("BUILDER_PARENT_DIR")
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	if runtime.GOOS == "windows" {
		artifactZip := artifactDir + ".zip"

		// CreateZip temp dir.
		outFile, err := os.Create(artifactZip)
		if err != nil {
			spinner.LogMessage("failed to create artifact directory: "+err.Error(), "fatal")
		}

		defer outFile.Close()

		// Create a new zip archive.
		w := zip.NewWriter(outFile)

		// Add files from artifact dir to the artifact zip.
		addFilesZip(w, artifactDir+"/", "")

		err = w.Close()
		if err != nil {
			spinner.LogMessage("failed to create artifact directory: "+err.Error(), "fatal")
		}
	} else {

		artifactTar := artifactDir + ".tar.gz"

		outFile, err := os.Create(artifactTar)
		if err != nil {
			spinner.LogMessage("failed to create artifact directory: "+err.Error(), "fatal")
		}

		defer outFile.Close()

		gw := gzip.NewWriter(outFile)
		defer gw.Close()
		tw := tar.NewWriter(gw)
		defer tw.Close()

		// Add files from artifact dir to the artifact tar.gz.
		addFilesTar(tw, artifactDir+"/", "")

		err = tw.Close()
		if err != nil {
			spinner.LogMessage("failed to create artifact: "+err.Error(), "fatal")
		}
	}

}

// recursively add files
func addFilesZip(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		spinner.LogMessage("failed to read zip directory: "+err.Error(), "fatal")
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
				spinner.LogMessage("failed to create zip: "+err.Error(), "error")
			}
			_, err = f.Write(dat)
			if err != nil {
				spinner.LogMessage("failed to add files to zip: "+err.Error(), "error")
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			addFilesZip(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func addFilesTar(w *tar.Writer, basePath, baseInZip string) {
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

			header, err := tar.FileInfoHeader(file, file.Name())
			if err != nil {
				fmt.Println(err)
			}

			header.Name = baseInZip + file.Name()

			// Add some files to the archive.
			err = w.WriteHeader(header)
			if err != nil {
				fmt.Println(err)
			}
			_, err = w.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			addFilesTar(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

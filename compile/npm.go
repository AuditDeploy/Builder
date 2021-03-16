package compile

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

//Npm creates zip from files passed in as arg
func Npm() {

	//copies contents of .hidden to workspace
	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	tempWorkspace := workspaceDir + "/temp/" 
	//make temp dir
	os.Mkdir(tempWorkspace, 0755)

	//add hidden dir contents to temp dir, install dependencies
	exec.Command("cp", "-a", hiddenDir+"/.", tempWorkspace).Run()
	exec.Command("npm", "install", tempWorkspace).Run()
	
	// Zip temp dir.
	outFile, err := os.Create(workspaceDir+"/temp.zip")
	if err != nil {
		 log.Fatal(err)
	}
	
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addFiles(w, tempWorkspace, "")

	err = w.Close()
	if err != nil {
		 log.Fatal(err)
	}

	//make hiddenDir hidden
	exec.Command("attrib", hiddenDir, "-h").Run()
	// Make contents read-only.
	exec.Command("chmod", "-R", "0444", hiddenDir).Run()
}

//recursively add files
func addFiles(w *zip.Writer, basePath, baseInZip string) {
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
					addFiles(w, newBase, baseInZip  + file.Name() + "/")
			}
	}
}
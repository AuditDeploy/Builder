package compile

import (
	"Builder/logger"
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Npm creates zip from files passed in as arg
func Npm() {

	hiddenDir := os.Getenv("BUILDER_HIDDEN_DIR")
	workspaceDir := os.Getenv("BUILDER_WORKSPACE_DIR")
	tempWorkspace := workspaceDir + "\\temp\\" 
	//make temp dir
	os.Mkdir(tempWorkspace, 0755)

	//add hidden dir contents to temp dir, install dependencies
	exec.Command("cp", "-a", hiddenDir+"/.", tempWorkspace).Run()

	//install dependencies/build, if yaml build type exists install accordingly
	buildTool := strings.ToLower(os.Getenv("BUILDER_BUILD_TOOL"))
	//get user working dir, add temp wrkSpace dir
	path, _ := os.Getwd()
	fullPath := path + tempWorkspace
	var cmd *exec.Cmd
	if (buildTool == "npm") {
		fmt.Println(buildTool)
		cmd = exec.Command("npm", "install") // or whatever the program is
    cmd.Dir = fullPath       // or whatever directory it's in
	} else {
		cmd = exec.Command("npm", "install") // or whatever the program is
    cmd.Dir = fullPath       // or whatever directory it's in
		//default
	}

	//run cmd, check for err, log cmd
	logger.InfoLogger.Println(cmd)
	err := cmd.Run()
	if err != nil {
		logger.ErrorLogger.Println("Node project failed to compile.")
		log.Fatal(err)
	} 

	// Zip temp dir.
	outFile, err := os.Create(workspaceDir+"/temp.zip")
	if err != nil {
		 log.Fatal(err)
	}
	
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files from temp dir to the archive.
	addNpmFiles(w, tempWorkspace, "")

	err = w.Close()
	if err != nil {
		logger.ErrorLogger.Println("Npm project failed to compile.")
		 log.Fatal(err)
	}

	artifactPath := os.Getenv("BUILDER_OUTPUT_PATH")
	if (artifactPath != "") {
		exec.Command("cp", "-a", workspaceDir+"/temp.zip", artifactPath).Run()
	}
	logger.InfoLogger.Println("Npm project compiled successfully.")
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
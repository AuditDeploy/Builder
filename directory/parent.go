package directory

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"Builder/spinner"
	"Builder/utils"
)

// MakeDirs does...
func MakeDirs() {
	//handles -n flag
	name := utils.GetName()

	//check for projectPath env from builder.yaml
	configPath := os.Getenv("BUILDER_DIR_PATH")

	// If file or folder already named 'builder' in path, change builder folder name to builder_data
	var builderFolderName string
	if configPath != "" {
		if _, err := os.Stat(configPath + "/" + "builder"); err != nil {
			builderFolderName = "builder"
		} else {
			builderFolderName = "builder_data"
		}
	} else {
		if _, err := os.Stat("./" + "builder"); err != nil {
			builderFolderName = "builder"
		} else {
			builderFolderName = "builder_data"
		}
	}

	var path string
	if os.Getenv("BUILDER_COMMAND") == "true" {
		if configPath != "" {
			if os.Getenv("BUILDER_BUILDS_DIR") != "" {
				path = configPath + "/" + name + "/" + os.Getenv("BUILDER_BUILDS_DIR") + "/" + name + "_" + name
			} else {
				path = configPath + "/" + name + "/builder/" + name + "_" + name
			}
		} else { // Place builds in builder folder in repo
			// Check if user wants to name builder folder a different name
			if os.Getenv("BUILDER_BUILDS_DIR") != "" {
				path = "./" + os.Getenv("BUILDER_BUILDS_DIR") + "/" + name + "_" + name
			} else {
				path = "./" + builderFolderName + "/" + name + "_" + name
			}
		}
	} else if os.Getenv("BUILDER_DOCKER_COMMAND") == "true" {
		if configPath != "" {
			path = configPath + "/" + name + "_" + name
		} else { // Place builds in builder_data folder in repo
			// Check if user wants to name builder folder a different name
			if os.Getenv("BUILDER_BUILDS_DIR") != "" {
				path = "./" + os.Getenv("BUILDER_BUILDS_DIR") + "/" + name + "_" + name
			} else {
				path = "./builder_data/" + name + "_" + name
			}
		}
	} else { // builder init so create an initial repo dir
		if configPath != "" {
			// Check if user wants to name builder folder a different name
			if os.Getenv("BUILDER_BUILDS_DIR") != "" {
				path = configPath + "/" + name + "/" + os.Getenv("BUILDER_BUILDS_DIR") + "/" + name + "_" + name
			} else {
				path = configPath + "/" + name + "/" + builderFolderName + "/" + name + "_" + name
			}
		} else {
			// Check if user wants to name builder folder a different name
			if os.Getenv("BUILDER_BUILDS_DIR") != "" {
				path = "./" + name + "/" + os.Getenv("BUILDER_BUILDS_DIR") + "/" + name + "_" + name
			} else {
				path = "./" + name + "/" + builderFolderName + "/" + name + "_" + name
			}
		}
	}

	if os.Getenv("BUILDER_DOCKER_COMMAND") == "true" {
		MakeParentDir(path)

		MakeWorkspaceDir(path)

		MakeLogsDir(path)
		MakeBuilderDir()
	} else {
		MakeParentDir(path)

		MakeHiddenDir(path)
		MakeWorkspaceDir(path)

		MakeLogsDir(path)
		MakeBuilderDir()
	}
}

func MakeParentDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	info, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
		spinner.LogMessage("Path already exists", "info")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) || !info.IsDir() {
		errDir := os.MkdirAll(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			spinner.LogMessage("failed to create directory at "+path+": "+err.Error(), "fatal")
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_PARENT_DIR")
	if !present {
		os.Setenv("BUILDER_PARENT_DIR", path)
	} else {
		fmt.Println("BUILDER_PARENT_DIR", val)
	}

	return true, err
}

func UpdateParentDirName(pathWithWrongParentName string) string {
	oldName, _ := os.LookupEnv("BUILDER_PARENT_DIR")
	projectName := utils.GetName()
	startTime, _ := time.Parse(time.RFC850, os.Getenv("BUILD_START_TIME"))
	unixTimestamp := startTime.Unix()
	newName := strings.TrimSuffix(oldName, projectName) + fmt.Sprint(unixTimestamp)

	path := os.Getenv("BUILDER_DIR_PATH")

	// user using builder command and did not provide a path to build in
	if os.Getenv("BUILDER_COMMAND") == "true" && path == "" {
		wdPath, err := os.Getwd()
		if err != nil {
			spinner.LogMessage("error getting builder working directory", "error")
		}

		path = wdPath
	}

	// user using builder docker command and did not provide a path to build in
	if os.Getenv("BUILDER_DOCKER_COMMAND") == "true" && path == "" {
		wdPath, err := os.Getwd()
		if err != nil {
			spinner.LogMessage("error getting builder working directory", "error")
		}

		path = wdPath
	}

	// Rename folder
	if oldName[0:2] == "./" {
		// err := os.Rename(path+"/"+oldName[2:], path+"/"+newName[2:])
		// if err != nil {
		// 	fmt.Println(err.(*os.LinkError).Err)
		// 	spinner.LogMessage("could not rename parent dir", "fatal")
		// }

		err := CopyDir(path+"/"+oldName[2:], path+"/"+newName[2:])
		if err != nil {
			spinner.LogMessage("Couldn't copy over files to new parent dir: "+err.Error(), "fatal")
		}

		err = os.RemoveAll(path + "/" + oldName[2:])
		if err != nil {
			spinner.LogMessage("Couldn't delete old parent dir: "+err.Error(), "fatal")
		}

	} else {
		// err := os.Rename(oldName, newName)
		// if err != nil {
		// 	fmt.Println(err.(*os.LinkError).Err)
		// 	spinner.LogMessage("could not rename parent dir", "fatal")
		// }

		err := CopyDir(oldName, newName)
		if err != nil {
			spinner.LogMessage("Couldn't copy over files to new parent dir: "+err.Error(), "fatal")
		}

		err = os.RemoveAll(oldName)
		if err != nil {
			spinner.LogMessage("Couldn't delete old parent dir: "+err.Error(), "fatal")
		}
	}

	// Update env vars to include new parent folder name
	name := utils.GetName()
	os.Setenv("BUILDER_PARENT_DIR", newName)
	os.Setenv("BUILDER_HIDDEN_DIR", newName+"/"+name)
	os.Setenv("BUILDER_WORKSPACE_DIR", newName+"/workspace")
	os.Setenv("BUILDER_LOGS_DIR", newName+"/logs")

	if os.Getenv("BUILDER_DOCKER_COMMAND") == "true" {
		oldArtifactPath := os.Getenv("BUILDER_ARTIFACT_DIR")
		newArtifactPath := strings.Replace(oldArtifactPath, name+"_"+name, name+"_"+fmt.Sprint(unixTimestamp), 1)

		os.Setenv("BUILDER_ARTIFACT_DIR", newArtifactPath)
	}

	// Return new path with new parent directory name
	newPath := strings.Replace(pathWithWrongParentName, oldName[2:], newName[2:], 1)

	return newPath
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

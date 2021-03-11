package directory

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func hiddenDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0444)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			log.Fatal(err)
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_HIDDEN_DIR")
	if !present {
		os.Setenv("BUILDER_HIDDEN_DIR", path)
		fmt.Println("BUILDER_HIDDEN_DIR", os.Getenv("BUILDER_HIDDEN_DIR"))
	} else {
		fmt.Println("BUILDER_HIDDEN_DIR", val)
	}

	//make directory hidden
	pathW, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fmt.Print(err)
	}
	err = syscall.SetFileAttributes(pathW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		fmt.Print(err)
	}

	return true, err
}

//MakeHiddenDir does...
func MakeHiddenDir(path string) {

	hiddenPath := path + "/.hidden"
	hiddenDir(hiddenPath)

}

package logs

import (
	"fmt"
	"log"
	"os"
)

func logDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {

		errDir := os.MkdirAll(path, 0755)
		//should return nil once directory is made, if not, throw err
		if errDir != nil {
			log.Fatal(err)
		}
		fmt.Println("Path created")

	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_LOGS_DIR")
	if !present {
		os.Setenv("BUILDER_LOGS_DIR", path)
		fmt.Println("BUILDER_LOGS_DIR", os.Getenv("BUILDER_LOGS_DIR"))
	} else {
		fmt.Println("BUILDER_LOGS_DIR", val)
	}

	//make directory hidden
	// pathW, err := syscall.UTF16PtrFromString(path)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// err = syscall.SetFileAttributes(pathW, syscall.FILE_ATTRIBUTE_HIDDEN)
	// if err != nil {
	// 	fmt.Print(err)
	// }

	return true, err
}

//MakeLogDir does...
func MakeLogDir(path string) {

	logPath := path + "/.logs"

	fmt.Printf(logPath)
	logDir(logPath)
}

package directory

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilarocca/Builder/utils"
	"github.com/manifoldco/promptui"
)

func parentDir(path string) (bool, error) {
	//check if file path exists, returns err = nil if file exists
	_, err := os.Stat(path)

	if err == nil {
		fmt.Println("Path already exists")
	}

	// should return true if file doesn't exist
	if os.IsNotExist(err) {
		//prompt user if they'd like dir to be created
		mk := yesNo()

		if mk == true {
			errDir := os.MkdirAll(path, 0755)
			//should return nil once directory is made, if not, throw err
			if errDir != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("Please create a directory for the Builder")
			os.Exit(1)
			return true, err
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

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Create A Directory? [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

//MakeParentDir does...
func MakeParentDir() {
	args := os.Args[1:]

	//handles -n flag
	name := utils.GetName(args)

	t := time.Now()

	// local path for now
	path := "./" + name +"_"+t.Format("20060102150405")

	parentDir(path)

	MakeHiddenDir(path)
	MakeLogDir(path)
	MakeWorkspaceDir(path)
}

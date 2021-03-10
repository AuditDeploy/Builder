package parent

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ilarocca/Builder/makeDirFunctions/sub/hidden"
	"github.com/ilarocca/Builder/makeDirFunctions/sub/logs"
	"github.com/ilarocca/Builder/makeDirFunctions/sub/workspace"

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
			fmt.Println("Path created")

		} else {
			fmt.Println("Please create a directory for the Builder")
			return true, err
		}
	}

	//check workspace env exists, if not, create it
	val, present := os.LookupEnv("BUILDER_PARENT_DIR")
	if !present {
		os.Setenv("BUILDER_PARENT_DIR", path)
		fmt.Println("BUILDER_PARENT_DIR", os.Getenv("BUILDER_PARENT_DIR"))
	} else {
		fmt.Println("BUILDER_PARENT_DIR", val)
	}
	return true, err
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Select[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

//MakeParentDir does...
func MakeParentDir(args string) {
	// 'github.com/name/project' slice that into '/name/project' as a var

	//slice original url
	name := args[strings.LastIndex(args, "/")+1:]

	// local path for now
	path := "./" + name


	fmt.Printf(path)
	parentDir(path)

	hidden.MakeHiddenDir(path)
	logs.MakeLogDir(path)
	workspace.MakeWorkspaceDir(path)
}

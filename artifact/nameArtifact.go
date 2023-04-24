package artifact

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//rename artifact with Unix timestamp
func NameArtifact(fullPath string, extName string) (string) {
	//seperate extName by last ".", return that ext (jar, exe, etc)
	newExtName := extName[strings.LastIndex(extName, ".")+1:]

	//trim off ".jar", ".exe", etc to add timestamp
	res := strings.Split(extName, "."+newExtName)
	currentTime := time.Now().Unix()

	//join it all back together
	artifactName := res[0]+"_"+strconv.FormatInt(currentTime, 10)+"."+newExtName

  err := os.Rename(fullPath+extName, fullPath+artifactName)

	if err != nil {
		fmt.Println(err)
	}

	return artifactName
}
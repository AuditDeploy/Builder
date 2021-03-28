package logger

import "strings"

//converts ./parentDirPath to fileName without "./"
func ConvertParentPathToFileName(filePath string) string {
	filePathArray := strings.Split(filePath, "")
	fileNameArray := filePathArray[2:]
	fileName := strings.Join(fileNameArray, "")
	return fileName
}

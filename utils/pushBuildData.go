package utils

import (
	"Builder/spinner"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func PushBuildData() {
	pushURL := os.Getenv("BUILDER_PUSH_URL")
	artifactDir := os.Getenv("BUILDER_ARTIFACT_DIR")

	// Read in build JSON data from build artifact directory
	metadataJSON, errMetadataRead := os.ReadFile(artifactDir + "/metadata.json")
	if errMetadataRead != nil {
		spinner.LogMessage("Can't find / open metadata file: "+errMetadataRead.Error(), "fatal")
	}

	// Read in logs JSON data from logs directory and add to metadata json
	lastSlash := strings.LastIndex(artifactDir, "/")
	logsFilePath := artifactDir[0:lastSlash] + "/logs/logs.json"
	validLogsJSONString := getLogsJSON(logsFilePath)

	var bodyFormat map[string]interface{}
	json.Unmarshal(metadataJSON, &bodyFormat)

	bodyFormat["logs"] = validLogsJSONString

	bodyJSON, errBodyMarshal := json.Marshal(bodyFormat)
	if errBodyMarshal != nil {
		spinner.LogMessage("Error encoding JSON body: "+errBodyMarshal.Error(), "fatal")
	}

	responseBody := bytes.NewBuffer(bodyJSON)

	// Send metadata+logs to user provided url
	req, _ := http.NewRequest("POST", pushURL, responseBody)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, errSend := client.Do(req)
	if errSend != nil {
		spinner.LogMessage("Error sending build metadata and logs: "+errSend.Error(), "fatal")
	} else {
		spinner.LogMessage("Metadata+logs data pushed to provided url", "info")
	}
	defer resp.Body.Close()

	// Send artifact(s) to user provided url
	artifactList := os.Getenv("BUILDER_ARTIFACT_LIST")
	artifactArray := strings.Split(artifactList, ",")

	buf := bytes.NewBuffer(nil)
	m := multipart.NewWriter(buf)

	for i, artifact := range artifactArray {
		part, err := m.CreateFormFile("artifact"+fmt.Sprint(i), artifact)
		if err != nil {
			spinner.LogMessage("Error posting artifact(s) to url.", "fatal")
		}

		file, err := os.Open(artifactDir + "/" + artifact)
		if err != nil {
			spinner.LogMessage("Error posting artifact(s) to url. Couldn't open artifact file: "+err.Error(), "fatal")
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}
	m.Close()

	artifactReq, _ := http.NewRequest("POST", pushURL, buf)
	req.Header.Add("Content-Type", m.FormDataContentType())

	artifactResp, errSend := client.Do(artifactReq)
	if errSend != nil {
		spinner.LogMessage("Error sending artifact(s) to push url: "+errSend.Error(), "fatal")
	}
	defer artifactResp.Body.Close()

	spinner.LogMessage("artifact(s) pushed to provided url", "info")
}

func getLogsJSON(path string) string {
	logsFile, err := os.Open(path)
	if err != nil {
		spinner.LogMessage("Can't find / open log file: "+err.Error(), "fatal")
	}
	defer logsFile.Close()

	var lines []string
	scanner := bufio.NewScanner(logsFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Re-format zap logger output to valid JSON
	asString := strings.Join(lines, ",")

	jsonString := "[" + asString + "]"

	return jsonString
}

package gui

import (
	"bufio"
	_ "embed"
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/zserge/lorca"
)

// Embed html code for gui
//
//go:embed gui_index.html
var IndexHtmlContents []byte

// Embed js code for gui
//
//go:embed gui.js
var jsContents []byte

// Embed css code for gui
//
//go:embed gui.css
var cssContents []byte

// Embed logo for gui
//
//go:embed logo.png
var logo []byte

type Build struct {
	ProjectName      string `json:"ProjectName"`
	ProjectType      string `json:"ProjectType"`
	ArtifactName     string `json:"ArtifactName"`
	ArtifactLocation string `json:"ArtifactLocation"`
	UserName         string `json:"UserName"`
	HomeDir          string `json:"HomeDir"`
	IP               string `json:"IP"`
	StartTime        string `json:"StartTime"`
	EndTime          string `json:"EndTime"`
	MasterGitHash    string `json:"MasterGitHash"`
	BranchName       string `json:"BranchName"`
	BuildHash        string `json:"BuildHash"`
}

func Gui() {

	// Read in json data
	getBuildsJSON := func() string {
		// Read in builds JSON data from application builder folder
		var buildsPath string
		if runtime.GOOS == "windows" {
			appDataDir := os.Getenv("LOCALAPPDATA")
			if appDataDir == "" {
				appDataDir = os.Getenv("APPDATA")
			}

			buildsPath = appDataDir + "/Builder/builds.json"
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}

			buildsPath = homeDir + "/.builder/builds.json"
		}

		buildsJSON, err := os.ReadFile(buildsPath)
		if err != nil {
			log.Fatal(err)
		}

		FormattedBuildsJSON := strings.TrimSuffix(string(buildsJSON), ",\n")

		FormattedBuildsJSON = "[" + FormattedBuildsJSON + "]"

		return FormattedBuildsJSON
	}

	getLogsJSON := func(path string) string {
		logsFile, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
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

	getImage := func() string {
		image := base64.StdEncoding.EncodeToString(logo)

		return image
	}

	// Combine html, css, and js files for gui
	cssRegex := regexp.MustCompile(`cssgoeshere`)
	jsRegex := regexp.MustCompile(`jsgoeshere`)

	finalHTMLContent := cssRegex.ReplaceAllString(string(IndexHtmlContents), string(cssContents))
	finalHTMLContent = jsRegex.ReplaceAllString(finalHTMLContent, string(jsContents))

	// Custom arguments for chrome popup
	args := []string{"--remote-allow-origins=*"}

	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("", "", 1200, 1000, args...)

	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("getBuildsJSON", getBuildsJSON)
	ui.Bind("getLogsJSON", getLogsJSON)
	ui.Bind("getImage", getImage)

	ui.Load("data:text/html," + url.PathEscape(finalHTMLContent))

	// Wait until UI window is closed
	<-ui.Done()
}

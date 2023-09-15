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

// Embed C logo
//
//go:embed CLogo.png
var CLogo []byte

// Embed CSharp logo
//
//go:embed CSharpLogo.png
var CSharpLogo []byte

// Embed Go logo
//
//go:embed GoLogo.png
var GoLogo []byte

// Embed Java logo
//
//go:embed JavaLogo.png
var JavaLogo []byte

// Embed Node logo
//
//go:embed NodeLogo.png
var NodeLogo []byte

// Embed Python logo
//
//go:embed PythonLogo.png
var PythonLogo []byte

// Embed Ruby logo
//
//go:embed RubyLogo.png
var RubyLogo []byte

// Embed Rust logo
//
//go:embed RustLogo.png
var RustLogo []byte

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

	getCLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(CLogo)

		return image
	}

	getCSharpLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(CSharpLogo)

		return image
	}

	getGoLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(GoLogo)

		return image
	}

	getJavaLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(JavaLogo)

		return image
	}

	getNodeLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(NodeLogo)

		return image
	}

	getPythonLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(PythonLogo)

		return image
	}

	getRubyLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(RubyLogo)

		return image
	}

	getRustLogoImage := func() string {
		image := base64.StdEncoding.EncodeToString(RustLogo)

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
	ui.Bind("getCLogoImage", getCLogoImage)
	ui.Bind("getCSharpLogoImage", getCSharpLogoImage)
	ui.Bind("getGoLogoImage", getGoLogoImage)
	ui.Bind("getJavaLogoImage", getJavaLogoImage)
	ui.Bind("getNodeLogoImage", getNodeLogoImage)
	ui.Bind("getPythonLogoImage", getPythonLogoImage)
	ui.Bind("getRubyLogoImage", getRubyLogoImage)
	ui.Bind("getRustLogoImage", getRustLogoImage)

	ui.Load("data:text/html," + url.PathEscape(finalHTMLContent))

	// Wait until UI window is closed
	<-ui.Done()
}

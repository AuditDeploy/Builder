package gui

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"time"

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
	Time        time.Time
	User        string
	Artifact    string
	ProjectName string
	GitHash     string
	BuildHash   string
}

func Gui() {

	// Read in json data
	getBuildsJSON := func() string {
		// Read in builds JSON data from .builder/ in user home dir
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		buildsJSON, err := os.ReadFile(homeDir + "/.builder/builds.json")
		if err != nil {
			log.Fatal(err)
		}

		return string(buildsJSON)
	}

	getLogsJSON := func(path string) string {
		logsJSON, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		return string(logsJSON)
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

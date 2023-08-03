package gui

import (
	_ "embed"
	"encoding/json"
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
var htmlContents []byte

// Embed js code for gui
//
//go:embed gui.js
var jsContents []byte

// Embed css code for gui
//
//go:embed gui.css
var cssContents []byte

type Build struct {
	Time        time.Time
	User        string
	Artifact    string
	ProjectName string
	GitHash     string
}

func Gui() {

	jsonToHTML := func() string {
		// Read in builds JSON data
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		buildsJSON, err := os.ReadFile(homeDir + "/.builder/builds.json")
		if err != nil {
			log.Fatal(err)
		}
		var builds []Build
		error := json.Unmarshal(buildsJSON, &builds)
		if error != nil {
			log.Fatal(error)
		}

		// Create HTML table
		text := ""
		for build := range builds {
			text += "</tr>"

			text += "<td>" + builds[build].Time.String() + "</td>"
			text += "<td>" + builds[build].User + "</td>"
			text += "<td>" + builds[build].Artifact + "</td>"
			text += "<td>" + builds[build].ProjectName + "</td>"
			text += "<td>" + builds[build].GitHash + "</td>"

			text += "</tr>"
		}

		return text
	}

	// Combine html, css, and js files for gui
	cssRegex := regexp.MustCompile(`cssgoeshere`)
	jsRegex := regexp.MustCompile(`jsgoeshere`)

	finalHTMLContent := cssRegex.ReplaceAllString(string(htmlContents), string(cssContents))
	finalHTMLContent = jsRegex.ReplaceAllString(finalHTMLContent, string(jsContents))

	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("", "", 1200, 1000, "--remote-allow-origins=*")

	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("jsonToHTML", jsonToHTML)

	ui.Load("data:text/html," + url.PathEscape(finalHTMLContent))

	// Wait until UI window is closed
	<-ui.Done()
}

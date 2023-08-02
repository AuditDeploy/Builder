package cmd

import (
	"log"
	"net/url"
	
	"github.com/zserge/lorca"
)

func Gui() {
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("data:text/html," + url.PathEscape(`
	<html>
		<head><title>Hello</title></head>
		<body><h1>Hellow, world!</h1></body>
	<html>
	`), "", 800, 800, "--remote-allow-origins=*")

	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()	

	// Wait until UI window is closed
	<-ui.Done()
}

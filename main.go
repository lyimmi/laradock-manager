package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func basic() string {
	cmd := exec.Command("docker-compose", "--version")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return "OK: " + out.String()
}

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    800,
		Title:     "Laradock manager",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})
	app.Bind(basic)
	app.Run()
}

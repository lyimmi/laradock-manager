package main

import (
	"bytes"
	"encoding/json"
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
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
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

	vuex := NewVuexState()
	res := VuexStore{}
	json.Unmarshal([]byte(vuex.Read()), &res)
	laradockPath := res.Settings["laradockPath"]

	dc := NewDockerCompose(laradockPath)

	app.Bind(dc)
	app.Bind(vuex)
	app.Run()
}

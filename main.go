package main

import (
	"encoding/json"
	"github.com/lyimmi/laradock-manager/docker"

	"github.com/leaanthony/mewn"

	"github.com/wailsapp/wails"
)

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
	err := json.Unmarshal([]byte(vuex.Read()), &res)
	if err != nil {
		panic(err)
	}
	laradockPath := res.Settings["laradockPath"]

	dc := docker.NewDockerCompose(laradockPath)

	app.Bind(dc)
	app.Bind(vuex)

	err = app.Run()

	if err != nil {
		panic(err)
	}
}

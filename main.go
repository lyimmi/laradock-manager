package main

import (
	"encoding/json"

	"github.com/leaanthony/mewn"
	"github.com/lyimmi/laradock-manager/docker"
	"github.com/lyimmi/laradock-manager/vuex"
	"github.com/wailsapp/wails"
)

// main
func main() {
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")
	projectJSON := mewn.String("./project.json")
	project := Project{}
	projectErr := json.Unmarshal([]byte(projectJSON), &project)
	if projectErr != nil {
		panic(projectErr)
	}
	app := wails.CreateApp(&wails.AppConfig{
		Width:     1280,
		Height:    800,
		Title:     "Laradock manager",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})

	vuexState := vuex.NewVuexState(project.Binaryname)
	app.Bind(&App{})
	app.Bind(docker.NewDockerCompose(vuexState))
	app.Bind(vuexState)

	err := app.Run()

	if err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"

	"github.com/leaanthony/mewn"
	"github.com/lyimmi/laradock-manager/docker"
	"github.com/lyimmi/laradock-manager/vuex"
	"github.com/wailsapp/wails"
)

// App wails runtime
type App struct {
	runtime *wails.Runtime
}

// WailsInit initialize wails
func (s *App) WailsInit(r *wails.Runtime) error {
	s.runtime = r
	return nil
}

// SelectDirectory open a directory selector dialog
func (s *App) SelectDirectory() string {
	return s.runtime.Dialog.SelectDirectory()
}

// SelectFile open a file selector dialog
func (s *App) SelectFile() string {
	return s.runtime.Dialog.SelectFile()
}

// Project represents project.json
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Version    string `json:"version"`
	Binaryname string `json:"binaryname"`
	Frontend   struct {
		Dir     string `json:"dir"`
		Install string `json:"install"`
		Build   string `json:"build"`
		Bridge  string `json:"bridge"`
		Serve   string `json:"serve"`
	} `json:"frontend"`
	WailsVersion string `json:"WailsVersion"`
}

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

	vuexState := vuex.NewVuexState(project.Name)
	dc := docker.NewDockerCompose(vuexState)
	myApp := &App{}

	app.Bind(myApp)
	app.Bind(dc)
	app.Bind(vuexState)

	err := app.Run()

	if err != nil {
		panic(err)
	}
}

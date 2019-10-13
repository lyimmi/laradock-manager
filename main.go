package main

import (
	"encoding/json"

	"github.com/leaanthony/mewn"
	"github.com/lyimmi/laradock-manager/docker"
	"github.com/wailsapp/wails"
)

// MyRuntime wails runtime
type MyRuntime struct {
	runtime *wails.Runtime
}

// WailsInit initialize wails
func (s *MyRuntime) WailsInit(r *wails.Runtime) error {
	s.runtime = r
	return nil
}

// SelectDirectory open a directory selector dialog
func (s *MyRuntime) SelectDirectory() string {
	return s.runtime.Dialog.SelectDirectory()
}

// SelectFile open a file selector dialog
func (s *MyRuntime) SelectFile() string {
	return s.runtime.Dialog.SelectFile()
}

// main
func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1280,
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
	dc := docker.NewDockerCompose(res.Settings["laradockPath"])
	my := &MyRuntime{}

	app.Bind(my)
	app.Bind(dc)
	app.Bind(vuex)

	err = app.Run()

	if err != nil {
		panic(err)
	}
}

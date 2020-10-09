package main

import "github.com/wailsapp/wails"

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

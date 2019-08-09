package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

func checkDotEnv() bool {
	if _, err := os.Stat("/home/lyimmi/Projects/webnetwork/laradock/.env"); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func copyEnv() bool {
	sourceFile, err := os.Open("/home/lyimmi/Projects/webnetwork/laradock/env-example")
	if err != nil {
		return false
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create("/home/lyimmi/Projects/webnetwork/laradock/.env")
	if err != nil {
		return false
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		return false
	}
	return bytesCopied > 0
}

func containers() string {
	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = "/home/lyimmi/Projects/webnetwork/laradock/"
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
	app.Bind(basic)
	app.Bind(containers)
	app.Bind(checkDotEnv)
	app.Bind(copyEnv)
	app.Run()
}

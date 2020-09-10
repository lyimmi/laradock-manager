package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

//Version app version
var Version string

//Project is the project json
var Project project

type project struct {
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

func main() {
	loadProjectFile()
	deleteBuildFolder()
	createBuildFolder()
	pack()
}

func getMainPath() string {
	v := getVersion()
	return path.Join("build", "laradock-manager-"+v)
}

func loadProjectFile() {
	jsonFile, err := os.Open("./project.json")
	defer jsonFile.Close()

	if err != nil {
		println("Probably not called this in a wrong working dir, try calling it in the laradock-manager's root!")
		println(err.Error())
	}
	j, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(j, &Project)
}

func getVersion() string {
	if Version != "" {
		return Version
	}

	Version = Project.Version
	return Version
}

func deleteBuildFolder() {
	err := os.RemoveAll(getMainPath())
	check(err)
}

func createBuildFolder() {
	writeControlFile()
	writeDesktop()
	copyIcon()
	copyBin()
}

func writeControlFile() {
	//Control file
	fp := path.Join(getMainPath(), "DEBIAN")
	err := os.MkdirAll(fp, os.ModePerm)
	check(err)

	cont :=
		`Package: ` + Project.Binaryname + `
Version: ` + getVersion() + `
Section: devel
Priority: optional
Architecture: amd64
Installed-Size: 4800
Maintainer: ` + Project.Author.Name + ` <` + Project.Author.Email + `>
Homepage: https://github.com/Lyimmi/laradock-manager
Description: ` + Project.Name + `
  A simple application for managing laradock containers.
  Developed and tested only on Ubuntu 20.04/19.04/18.04
 .
  Made with https://wails.app/ (go & vue.js & vuetify)
 .
  Usage
  In order to use this your current user need to be able to access docker without sudo
  Create the docker group: $ sudo groupadd docker
  Add your user to the docker group: $ sudo usermod -aG docker $USER
  Log out and log back in so that your group membership is re-evaluated. ($ newgrp docker)
`

	f, err := os.Create(path.Join(fp, "control"))
	check(err)
	_, err = f.WriteString(cont)
	check(err)
	f.Close()
}

func writeDesktop() {
	fp := path.Join(getMainPath(), "usr", "share", "applications")
	err := os.MkdirAll(fp, os.ModePerm)
	check(err)

	cont := `[Desktop Entry]
Version=1.1
Type=Application
Name=Laradock Manager
Comment=Simple application for managing laradock containers
Icon=/opt/laradock-manager/laradock-manager.png
Exec=/opt/laradock-manager/laradock-manager
Actions=
Categories=Development;
StartupNotify=true`

	f, err := os.Create(path.Join(fp, "laradock-manager.desktop"))
	check(err)
	_, err = f.WriteString(cont)
	check(err)
	f.Close()
}

func copyIcon() {
	var err error
	//Menu icon
	fp := path.Join(getMainPath(), "opt", "laradock-manager")
	err = os.MkdirAll(fp, os.ModePerm)
	check(err)
	src := "assets/laradock-manager.png"
	dst := path.Join(fp, "laradock-manager.png")
	_, err = copy(src, dst)
	check(err)
}

func copyBin() {
	var err error
	fp := path.Join(getMainPath(), "opt", "laradock-manager")
	err = os.MkdirAll(fp, os.ModePerm)
	check(err)
	dst := path.Join(fp, "laradock-manager")
	_, err = copy("build/laradock-manager", dst)

	c := exec.Command("chmod", "+x", dst)
	err = c.Run()
	check(err)

}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func pack() {
	var err error
	var out []byte

	cmd := exec.Command("dpkg-deb", "--build", "build/laradock-manager-"+getVersion())
	// cmd.Dir = d
	out, err = cmd.CombinedOutput()
	fmt.Println(string(out))
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

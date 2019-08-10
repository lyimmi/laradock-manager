package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//DockerCompose DockerCompose struct
type DockerCompose struct {
	laradockPath        string
	availableContainers map[string]string
}

//NewDockerCompose Create a new DockerCompose struct
func NewDockerCompose(path string) *DockerCompose {
	result := &DockerCompose{laradockPath: path}
	return result
}

//SetLaradockPath Check if .env file exists
func (t *DockerCompose) SetLaradockPath(path string) bool {
	t.laradockPath = path
	return true
}

//CheckDotEnv Check if .env file exists
func (t *DockerCompose) CheckDotEnv() bool {
	if _, err := os.Stat(filepath.Join(t.laradockPath, ".env")); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

//CopyEnv Make the .env file form env-example
func (t *DockerCompose) CopyEnv() bool {
	sourceFile, err := os.Open(filepath.Join(t.laradockPath, "env-example"))
	if err != nil {
		return false
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(filepath.Join(t.laradockPath, ".env"))
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

//GetContainers run docker-compose ps and parse the output
func (t *DockerCompose) GetContainers() string {
	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = filepath.Join(t.laradockPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		//fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}

	reg := regexp.MustCompile(`\n`)
	lines := reg.Split(out.String(), -1)
	var c [][]string
	for _, e := range lines {
		reg = regexp.MustCompile(`\s\s+`)
		c = append(c, reg.Split(e, -1))
	}

	s, err := json.Marshal(c)
	if err != nil {
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return string(s)
}

//GetAvailableContainers run docker-compose ps and parse the output
func (t *DockerCompose) GetAvailableContainers() string {
	cmd := exec.Command("docker-compose", "config")
	cmd.Dir = filepath.Join(t.laradockPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}

	reg := regexp.MustCompile(`\n`)
	lines := reg.Split(out.String(), -1)
	var c []string
	gather := false
	for _, e := range lines {
		if !gather && strings.Contains(e, "services:") {
			gather = true
		}
		if strings.HasPrefix(e, "  ") && !strings.HasPrefix(e, "   ") && gather {
			e = strings.TrimSuffix(e, ":")
			c = append(c, strings.TrimSpace(e))
			fmt.Println(e)
		}
		if strings.Contains(e, "version:") {
			gather = false
			break
		}

	}

	s, err := json.Marshal(c)
	if err != nil {
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return string(s)
}

//ToggleContainer Toggle a container on and off
func (t *DockerCompose) ToggleContainer(state string, container string) bool {
	cmd := exec.Command("docker-compose", state, container)
	cmd.Dir = t.laradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return false
	}
	fmt.Println(out.String())
	return true
}

func (t *DockerCompose) regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

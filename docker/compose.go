package docker

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
	"time"
)

//Compose DockerCompose struct
type Compose struct {
	laradockPath        string
	containerExec       bool
	containerConnected  bool
	availableContainers map[string]string
}

// //WailsInit wails init
// func (t *Compose) WailsInit(runtime *wails.Runtime) error {
// 	go func() {
// 		for {
// 			if t.containerExec == true && t.containerConnected == true {
// 				runtime.Events.Emit("containerExecOutputChange", "connected")
// 				fmt.Println("event emitted")
// 			} else if t.containerExec == false && t.containerConnected == true {
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// 	return nil
// }

//NewDockerCompose Create a new DockerCompose struct
func NewDockerCompose(path string) *Compose {
	result := &Compose{laradockPath: path}
	return result
}

//SetLaradockPath Check if .env file exists
func (t *Compose) SetLaradockPath(path string) bool {
	t.laradockPath = path
	return true
}

//CheckDotEnv Check if .env file exists
func (t *Compose) CheckDotEnv() bool {
	if _, err := os.Stat(filepath.Join(t.laradockPath, ".env")); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

//CopyEnv Make the .env file form env-example
func (t *Compose) CopyEnv() bool {
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
func (t *Compose) GetContainers() string {
	cmd := exec.Command("docker-compose", "ps")
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
func (t *Compose) GetAvailableContainers() string {
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
			//fmt.Println(e)
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
func (t *Compose) ToggleContainer(state string, container string) bool {
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
	return true
}

//ExecContainer Execute a docker container
func (t *Compose) ExecContainer(container string, user string) {
	t.containerExec = true
	fmt.Println("connecting composer exec")
	stopExecChanel := make(chan bool)
	go func() {
		containerExec(user, container, t.laradockPath, stopExecChanel)
	}()
	for {
		if t.containerExec == false {
			fmt.Println("disconnecting exec")
			close(stopExecChanel)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

//StopExecContainer Stop the exec's g routine
func (t *Compose) StopExecContainer() string {
	t.containerExec = false
	fmt.Println("disconnect signal reveived")
	return "disconnected"
}

func (t *Compose) regSplit(text string, delimeter string) []string {
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

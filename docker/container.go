package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/lyimmi/laradock-manager/vuex"
	"github.com/wailsapp/wails"
)

//Compose DockerCompose struct
type Compose struct {
	vuexState    *vuex.State
	runtime      *wails.Runtime
	statsQuit    chan struct{}
	statsRunning bool
	HasUp        bool
	HasUpLast    time.Time
}

type container struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Favorite bool   `json:"favorite"`
}

// WailsInit init
func (t *Compose) WailsInit(runtime *wails.Runtime) error {
	t.runtime = runtime
	return nil
}

//NewDockerCompose Create a new DockerCompose struct
func NewDockerCompose(vuexState *vuex.State) *Compose {
	compose := &Compose{vuexState: vuexState}
	vuexState.Read()
	return compose
}

//SetLaradockPath Check if .env file exists
func (t *Compose) SetLaradockPath(path string) bool {
	t.vuexState.Store.Settings.LaradockPath = path
	return true
}

//SetTerminalPath Check if .env file exists
func (t *Compose) SetTerminalPath(path string) bool {
	t.vuexState.Store.Settings.TerminalPath = path
	return true
}

//CheckDockerComposeVersion Check the docker-compose executable's version
func (t *Compose) CheckDockerComposeVersion() string {
	cmd := exec.Command("docker-compose", "-v")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return returnResponse(false, fmt.Sprint(err)+": "+stderr.String())
	}
	return returnResponse(true, out.String())
}

//GetContainersWithStatusesSlice get a slice of statuses
func (t *Compose) GetContainersWithStatusesSlice() ([]container, error) {
	var stdout, stderr bytes.Buffer
	var containers []container

	//Get docker-compose ps
	cmd := exec.Command("docker-compose", "--no-ansi", "ps", "-a")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := string(stdout.Bytes())
	if err != nil {
		return make([]container, 0), err
	}
	reg := regexp.MustCompile("\n" + t.vuexState.Store.Settings.ContainerPrefix + "_")
	lines := reg.Split(outStr, -1)
	for i, e := range lines {
		reg = regexp.MustCompile(`\s\s+`)
		contArr := reg.Split(e, -1)
		if i > 0 {
			if len(contArr) > 2 {
				name := strings.Replace(contArr[0], "_1", "", -1)
				if name == "docker-in-" { //fix docker-in-docker container...
					name = name + "docker"
				}
				state := "Up"
				favorite := false
				for _, f := range t.vuexState.Store.Containers.Favorites {
					if f == name {
						favorite = true
					}
				}
				if strings.Contains(contArr[2], "Exit") {
					state = "Stopped"
				}
				cont := container{
					Name:     name,
					State:    state,
					Favorite: favorite,
				}
				containers = append(containers, cont)
			}
		}
	}

	//Get all available containers from yml containers
	stdout, stderr = bytes.Buffer{}, bytes.Buffer{}
	outStr = ""

	cmd = exec.Command("docker-compose", "--no-ansi", "ps", "--services")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr = string(stdout.Bytes())
	if err != nil {
		return make([]container, 0), err
	}

	reg = regexp.MustCompile(`\n`)
	lines = reg.Split(outStr, -1)
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		ok := true
		for _, c := range containers {
			if c.Name == line {
				ok = false
				break
			}
		}
		if ok {
			favorite := false
			for _, f := range t.vuexState.Store.Containers.Favorites {
				if f == line {
					favorite = true
				}
			}
			cont := container{
				Name:     line,
				State:    "Down",
				Favorite: favorite,
			}
			containers = append(containers, cont)
		}
	}
	return containers, nil
}

//GetContainersWithStatuses run docker-compose ps and parse the output
func (t *Compose) GetContainersWithStatuses() string {
	containers, err := t.GetContainersWithStatusesSlice()
	if err != nil {
		t.emitError("Error: " + fmt.Sprint(err))
		return returnResponse(false, "Error: "+fmt.Sprint(err))
	}
	res, resErr := json.Marshal(containers)
	if resErr != nil {
		t.emitError("Error: " + fmt.Sprint(resErr))
		return returnResponse(false, "Error: "+fmt.Sprint(resErr))
	}
	return returnResponse(true, string(res))
}

//HasRunning check if there are any containers running
func (t *Compose) HasRunning() bool {
	if time.Now().After(t.HasUpLast.Add(2 * time.Second)) {
		containers, err := t.GetContainersWithStatusesSlice()
		if err != nil {
			t.emitError("Error: " + fmt.Sprint(err))
			return false
		}
		ok := false
		for _, c := range containers {
			if c.State == "Up" {
				ok = true
				break
			}
		}
		t.HasUp = ok
		return t.HasUp
	}
	return t.HasUp

}

//GetContainers run docker-compose ps --services and parse the output
func (t *Compose) GetContainers() string {
	cmd := exec.Command("docker-compose", "ps", "--services")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, fmt.Sprint(err)+": "+stderr.String())
	}

	reg := regexp.MustCompile(`\n`)
	lines := reg.Split(out.String(), -1)
	s, err := json.Marshal(lines)
	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, fmt.Sprint(err)+": "+stderr.String())
	}
	return returnResponse(true, string(s))
}

//Toggle Toggle a container on and off
func (t *Compose) Toggle(state string, containers string) string {
	cSlice := strings.Split(containers, "|")       //split the provided containers into a slice
	args := []string{state}                        //prepare args
	args = append(args, cSlice...)                 //merge all arguments
	cmd := exec.Command("docker-compose", args...) //build command
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, "false")
	}
	return returnResponse(true, "true")
}

//Down Down all the containers
func (t *Compose) Down() string {
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, "false")
	}
	return returnResponse(true, "true")
}

//Up Up a container
func (t *Compose) Up(containers string) string {
	cSlice := strings.Split(containers, "|")   //split the provided containers into a slice
	args := []string{"up", "--no-build", "-d"} //prepare args
	args = append(args, cSlice...)             //merge all arguments
	fmt.Println(args)
	cmd := exec.Command("docker-compose", args...) //build command
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, "false")
	}
	return returnResponse(true, "true")
}

//Build Build a container
func (t *Compose) Build(containers string, force bool) string {
	cSlice := strings.Split(containers, "|") //split the provided containers into a slice
	//prepare args
	var args []string
	if force == true {
		args = []string{"build", "--no-cache"}
	} else {
		args = []string{"build"}
	}
	args = append(args, cSlice...)                 //merge all arguments
	cmd := exec.Command("docker-compose", args...) //build command

	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		t.emitError(fmt.Sprint(err) + ": " + stderr.String())
		return returnResponse(false, "false")
	}
	return returnResponse(true, "true")
}

//Exec Execute a docker container
//
//start cmd /k echo Hello, World!
func (t *Compose) Exec(container string, user string) string {
	cmd := exec.Command("gnome-terminal", "--", "docker-compose", "exec", "--user="+user, container, "bash")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", "cmd", "/k", "docker-compose", "exec", "--user="+user, container, "bash")
	}
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	if err := cmd.Run(); err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}
	return returnResponse(true, "terminal started")
}

// Logs show logs
func (t *Compose) Logs(container string) string {
	cmd := exec.Command("gnome-terminal", "--", "docker-compose", "logs", "-f", "--tail=100", container)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", "cmd", "/k", "docker-compose", "logs", "-f", "--tail=100", container)
	}
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	if err := cmd.Run(); err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}
	return returnResponse(true, "Logs Started")
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

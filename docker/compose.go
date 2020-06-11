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
	"runtime"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/lyimmi/laradock-manager/vuex"
	"github.com/wailsapp/wails"
)

//Compose DockerCompose struct
type Compose struct {
	vuexState    *vuex.State
	runtime      *wails.Runtime
	statsQuit    chan struct{}
	statsRunning bool
}

// envStruct
type envStruck struct {
}

// Response struct
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type container struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Favorite bool   `json:"favorite"`
}

type stats struct {
	Name                string `json:"name"`
	CPUPercString       string `json:"cpu_perc_string"`
	CPUPerc             string `json:"cpu_perc"`
	MemoryUseage        string `json:"memory_usage"`
	MemoryPercentString string `json:"memory_percent_string"`
	MemoryPercent       string `json:"memory_percent"`
}

// returnResponse a response json
func returnResponse(success bool, message string) string {
	resp, err := json.Marshal(Response{Success: success, Message: message})
	if err != nil {
		log.Error(err)
	}
	return string(resp)
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

//CheckDotEnv Check if .env file exists
func (t *Compose) CheckDotEnv() string {
	if _, err := os.Stat(filepath.Join(t.vuexState.Store.Settings.LaradockPath, ".env")); err != nil {
		return returnResponse(true, "false")
	}
	return returnResponse(true, "true")
}

//CheckDockerVersion Check the docker executable's version
func (t *Compose) CheckDockerVersion() string {
	cmd := exec.Command("docker", "-v")
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

//CopyEnv Make the .env file form env-example
func (t *Compose) CopyEnv() string {
	sourceFile, err := os.Open(filepath.Join(t.vuexState.Store.Settings.LaradockPath, "env-example"))
	if err != nil {
		return returnResponse(false, "false")
	}

	// Create new file
	newFile, err := os.Create(filepath.Join(t.vuexState.Store.Settings.LaradockPath, ".env"))
	if err != nil {
		return returnResponse(false, "false")
	}

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		return returnResponse(false, "false")
	}

	sourceFile.Close()
	newFile.Close()

	if bytesCopied > 0 {
		return returnResponse(false, "true")
	}
	return returnResponse(false, "false")
}

func (t *Compose) listContainers() {
}

//GetContainersWithStatuses run docker-compose ps and parse the output
func (t *Compose) GetContainersWithStatuses() string {
	var stdout, stderr bytes.Buffer
	var containers []container

	//Get docker-compose ps
	cmd := exec.Command("docker-compose", "--no-ansi", "ps", "-a")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		t.emitError(errStr)
		return returnResponse(false, "Error: "+fmt.Sprint(err)+": "+errStr)
	}
	reg := regexp.MustCompile("\n" + t.vuexState.Store.Settings.ContainerPrefix + "_")
	lines := reg.Split(outStr, -1)
	for i, e := range lines {
		reg = regexp.MustCompile(`\s\s+`)
		contArr := reg.Split(e, -1)
		if i > 0 {
			if len(contArr) > 2 {
				name := strings.Replace(contArr[0], "_1", "", -1)
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
	outStr, errStr = "", ""

	cmd = exec.Command("docker-compose", "--no-ansi", "ps", "--services")
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		t.emitError("Error: " + fmt.Sprint(err) + ": " + errStr)
		return returnResponse(false, "Error: "+fmt.Sprint(err)+": "+errStr)
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
	res, resErr := json.Marshal(containers)
	if resErr != nil {
		t.emitError("Error: " + fmt.Sprint(resErr))
		return returnResponse(false, "Error: "+fmt.Sprint(resErr))
	}
	return returnResponse(true, string(res))
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
	cSlice := strings.Split(containers, "|")       //split the provided containers into a slice
	args := []string{"up", "-d", "--no-build"}     //prepare args
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
	cmd := exec.Command("gnome-terminal", "--", "docker-compose", "logs", "-f", "--tail=100")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", "cmd", "/k", "docker-compose", "logs", "-f", "--tail=100")
	}
	cmd.Dir = t.vuexState.Store.Settings.LaradockPath
	if err := cmd.Run(); err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}
	return returnResponse(true, "Logs Started")
}

//DotEnvContent Return dot env contents
func (t *Compose) DotEnvContent() map[string]string {
	env, err := godotenv.Read(filepath.Join(t.vuexState.Store.Settings.LaradockPath, ".env"))
	if err != nil {
		t.emitError("Error loading .env file")
		log.Fatal(err)
	}
	return env
}

//SaveDotEnvContent save dot env contents
func (t *Compose) SaveDotEnvContent(data string) string {
	f, err := os.Create(filepath.Join(t.vuexState.Store.Settings.LaradockPath, ".env"))
	if err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}

	_, err = f.WriteString(data)
	if err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}

	err = f.Sync()
	if err != nil {
		return returnResponse(false, fmt.Sprint(err))
	}

	f.Close()
	return returnResponse(true, "true")
}

// Stats collects docker stats
func (t *Compose) Stats() {
	if !t.statsRunning {
		t.statsRunning = true
		t.statsQuit = make(chan struct{})
		go func() {
			running := false
			for {
				select {
				case <-t.statsQuit:
					return
				default:
					if !running {
						running = true
						cmd := exec.Command("docker", "stats", "--no-stream", "--format", "\"{{.Name}}\\t{{.CPUPerc}}\\t{{.MemUsage}}\\t{{.MemPerc}}\"")
						out, err := cmd.Output()
						if err != nil {
							t.emitError(err.Error())
						}
						reg := regexp.MustCompile("\n")
						lines := reg.Split(string(out), -1)
						lines = lines[:len(lines)-1]
						statsa := []stats{}
						for _, line := range lines {
							reg = regexp.MustCompile("\t")
							contArr := reg.Split(line, -1)
							name := strings.Replace(contArr[0], "_1", "", -1)
							name = strings.Replace(name, t.vuexState.Store.Settings.ContainerPrefix+"_", "", -1)
							name = strings.Replace(name, `"`, "", -1)
							stat := stats{
								Name:                name,
								CPUPercString:       contArr[1],
								CPUPerc:             strings.Replace(contArr[1], `%`, "", -1),
								MemoryUseage:        contArr[2],
								MemoryPercentString: strings.Replace(contArr[3], `"`, "", -1),
								MemoryPercent:       strings.Replace(contArr[3], `%"`, "", -1),
							}
							statsa = append(statsa, stat)
						}
						res, mErr := json.Marshal(statsa)
						if mErr != nil {
							t.emitError(err.Error())
						}
						uEnc := b64.URLEncoding.EncodeToString(res)
						t.runtime.Events.Emit("stats", uEnc)
						running = false
					}
				}
				time.Sleep(5 * time.Second)
			}
		}()
	}
}

//StatsStop stops docker stats go routine
func (t *Compose) StatsStop() {
	if t.statsRunning {
		close(t.statsQuit)
	} else {
		t.emitError("statse allready closed")
	}
	t.statsRunning = false
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

func (t *Compose) emitError(err string) {
	uEnc := b64.URLEncoding.EncodeToString([]byte(err))
	t.runtime.Events.Emit("backendError", uEnc)
}

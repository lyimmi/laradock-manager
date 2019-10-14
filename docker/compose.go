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

	"github.com/joho/godotenv"
)

//Compose DockerCompose struct
type Compose struct {
	laradockPath        string
	terminalPath        string
	containerExec       bool
	containerConnected  bool
	availableContainers map[string]string
}

// envStruct
type envStruck struct {
}

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

//SetTerminalPath Check if .env file exists
func (t *Compose) SetTerminalPath(path string) bool {
	t.terminalPath = path
	return true
}

//CheckDotEnv Check if .env file exists
func (t *Compose) CheckDotEnv() bool {
	if _, err := os.Stat(filepath.Join(t.laradockPath, ".env")); err != nil {
		return !os.IsNotExist(err)
	}
	return true
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
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return out.String()
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
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return out.String()
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

//Get run docker-compose ps and parse the output
func (t *Compose) Get() string {
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

//GetAvailables run docker-compose ps --services and parse the output
func (t *Compose) GetAvailables() string {
	cmd := exec.Command("docker-compose", "ps", "--services")
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
	s, err := json.Marshal(lines)
	if err != nil {
		return "Error: " + fmt.Sprint(err) + ": " + stderr.String()
	}
	return string(s)
}

//Toggle Toggle a container on and off
func (t *Compose) Toggle(state string, container string) bool {
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

//Down Down all the containers
func (t *Compose) Down() bool {
	cmd := exec.Command("docker-compose", "down")
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

//Up Up a container
func (t *Compose) Up(containers string) bool {
	cSlice := strings.Split(containers, "|")       //split the provided containers into a slice
	args := []string{"up", "-d", "--no-build"}     //prepare args
	args = append(args, cSlice...)                 //merge all arguments
	cmd := exec.Command("docker-compose", args...) //build command
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

//Build Build a container
func (t *Compose) Build(containers string, force bool) bool {
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

//Exec Execute a docker container
//
//start cmd /k echo Hello, World!
func (t *Compose) Exec(container string, user string) string {
	cmd := exec.Command("gnome-terminal", "--", "docker-compose", "exec", "--user="+user, container, "bash")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", "cmd", "/k", "docker-compose", "exec", "--user="+user, container, "bash")
	}
	cmd.Dir = filepath.Join(t.laradockPath)
	if err := cmd.Run(); err != nil {
		return "Error: " + fmt.Sprint(err)
	}
	return "terminal started"
}

//DotEnvContent Return dot env contents
func (t *Compose) DotEnvContent() map[string]string {
	env, err := godotenv.Read(filepath.Join(t.laradockPath, ".env"))
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return env
}

//SaveDotEnvContent save dot env contents
func (t *Compose) SaveDotEnvContent(data string) bool {
	f, err := os.Create(filepath.Join(t.laradockPath, ".env"))
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return false
	}
	defer f.Close()

	n, err := f.WriteString(data)
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return false
	}
	fmt.Printf("wrote %d bytes\n", n)

	err = f.Sync()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return false
	}

	return true
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

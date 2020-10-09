package docker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// envStruct
type envStruck struct {
}

//CheckDotEnv Check if .env file exists
func (t *Compose) CheckDotEnv() bool {
	if _, err := os.Stat(filepath.Join(t.vuexState.Store.Settings.LaradockPath, ".env")); err != nil {
		return false
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

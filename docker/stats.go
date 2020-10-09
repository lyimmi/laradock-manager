package docker

import (
	b64 "encoding/base64"
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type stats struct {
	Name                string `json:"name"`
	CPUPercString       string `json:"cpu_perc_string"`
	CPUPerc             string `json:"cpu_perc"`
	MemoryUseage        string `json:"memory_usage"`
	MemoryPercentString string `json:"memory_percent_string"`
	MemoryPercent       string `json:"memory_percent"`
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
	}
	t.statsRunning = false
}

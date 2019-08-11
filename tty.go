package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/creack/pty"
)

func main() {
	fmt.Println("ExecContainer")
	go func() {
		fmt.Println("in process")
		cmd := exec.Command("docker-compose", "exec", "--user=laradock", "workspace", "bash")
		cmd.Dir = filepath.Join("/home/lyimmi/Projects/webnetwork/laradock")
		f, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}

		for {

			// io.Copy(os.Stdout, f)
			fmt.Println(f)
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		fmt.Println("aaa")
		time.Sleep(10 * time.Second)
	}
}

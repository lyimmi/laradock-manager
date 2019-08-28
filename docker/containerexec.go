package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kr/pty"
	"golang.org/x/net/websocket"
)

const (
	webAssetsDir = "assets/exec-www"
	listenAddr   = "127.0.0.1:5000"
)

//Handler Handler struckt
type Handler struct {
	fileServer http.Handler
}

var user string
var container string
var laradockPath string
var chanel chan bool
var isRunning = false

func containerExec(u string, c string, l string, ch chan bool) *http.Server {
	user = u
	container = c
	laradockPath = l
	chanel = ch

	//Ugly as hell hack close your eyes....

	srv := &http.Server{Addr: listenAddr}
	fmt.Printf("Listening on http://%s\n", listenAddr)
	handler := Handler{
		fileServer: http.FileServer(http.Dir(webAssetsDir)),
	}
	srv.Handler = &handler
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			isRunning = false
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v %v", r.Method, r.URL.Path, user)

	// need to serve shell via websocket?
	if strings.Trim(r.URL.Path, "/") == "shell" {
		onShell(w, r, user, container, laradockPath, chanel)
		return
	}
	// serve static assets from 'static' dir:
	h.fileServer.ServeHTTP(w, r)
}

// GET /shell handler
// Launches /bin/bash and starts serving it via the terminal
func onShell(w http.ResponseWriter, r *http.Request, user string, container string, lpath string, ch chan bool) {
	wsHandler := func(ws *websocket.Conn) {

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			// wrap the websocket into UTF-8 wrappers:
			wrapper := NewWebSockWrapper(ws, WebSocketTextMode)
			stdout := wrapper
			stderr := wrapper

			// this one is optional (solves some weird issues with vim running under shell)
			stdin := &InputWrapper{ws}

			cmd := exec.CommandContext(ctx, "docker-compose", "exec", container, "bash")
			if user != "" {
				cmd = exec.CommandContext(ctx, "docker-compose", "exec", "--user="+user, container, "bash")
			}
			cmd.Dir = filepath.Join(lpath)
			tty, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			defer tty.Close()

			// pipe to/fro websocket to the TTY:
			go func() {
				io.Copy(stdout, tty)
			}()
			go func() {
				io.Copy(stderr, tty)
			}()
			go func() {
				io.Copy(tty, stdin)
			}()
			// wait for the command to exit, then close the websocket
			cmd.Wait()
		}()
		for {
			select {
			default:
				time.Sleep(1 * time.Second)
			case <-ch:
				cancel()
				return
				// stop
			}
		}
	}
	defer log.Printf("Websocket session closed for %v", r.RemoteAddr)

	// start the websocket session:
	websocket.Handler(wsHandler).ServeHTTP(w, r)

}

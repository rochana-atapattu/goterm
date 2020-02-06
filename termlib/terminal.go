package terminal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/platform"
	"github.com/liamg/aminal/terminal"
	"github.com/riywo/loginshell"
)

type instance struct {
	Instanceid string `json:"instanceid"`
}
type body struct {
	Body   obj `json:"body"`
	Status int `json:"statusCode"`
}
type bastion struct {
	Publicip    string `json: "publicip"`
	BastionUser string `json: "bastionuser"`
}
type obj struct {
	Privateip string    `json: "privateip"`
	User      string    `json: "user"`
	Bastions  []bastion `json: "bastions"`
}

func CreateTerm(cmd string) platform.Pty {
	conf := getConfig()
	logger, err := getLogger(conf)
	if err != nil {
		fmt.Printf("Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Infof("Allocating pty...")

	pty, err := platform.NewPty(80, 25)
	if err != nil {
		logger.Fatalf("Failed to allocate pty: %s", err)
	}

	shellStr := conf.Shell
	if shellStr == "" {
		loginShell, err := loginshell.Shell()
		if err != nil {
			logger.Fatalf("Failed to ascertain your shell: %s", err)
		}
		shellStr = loginShell
	}

	os.Setenv("TERM", "xterm-256color") // controversial! easier than installing terminfo everywhere, but obviously going to be slightly different to xterm functionality, so we'll see...
	os.Setenv("COLORTERM", "truecolor")

	guestProcess, err := pty.CreateGuestProcess(shellStr)
	if err != nil {
		pty.Close()
		logger.Fatalf("Failed to start your shell: %s", err)
	}
	defer guestProcess.Close()

	pty.Write([]byte(cmd))

	logger.Infof("Creating terminal...")
	terminal := terminal.New(pty, logger, conf)

	g, err := gui.New(conf, terminal, logger)
	if err != nil {
		logger.Fatalf("Cannot start: %s", err)
	}

	go func() {
		if err := guestProcess.Wait(); err != nil {
			logger.Fatalf("Failed to wait for guest process: %s", err)
		}
		g.Close()
	}()

	fmt.Println("got request guestProcess.Wait()")
	func() {
		runtime.LockOSThread()
		if err := g.Render(); err != nil {
			logger.Fatalf("Render error: %s", err)
		}
	}()
	logger.Infof("Done . . .")

	return pty
}

func GetServers(instanceID string) body {
	jsonData := map[string]string{"instance_id": instanceID}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("https://jj75ccpa08.execute-api.ap-southeast-1.amazonaws.com/dev/server", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	println(string(data))

	b := body{}
	json.Unmarshal(data, &b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	return b
}

func CreateProxyCmd(b body) string {
	proxycommand := fmt.Sprintf("ssh -J %s@%s:22 %s@%s \r", b.Body.Bastions[0].BastionUser, b.Body.Bastions[0].Publicip, b.Body.User, b.Body.Privateip)
	return proxycommand
}

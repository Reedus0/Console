package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	ic "github.com/WAY29/icecream-go/icecream"
)

var (
	Bots = map[string]*Bot{}
)

type Bot struct {
	Name      string
	Logs      [10]string
	Timestamp int64
}

func (s *Bot) Json() string {
	bts, _ := json.Marshal(s)
	return string(bts)

}
func (s *Bot) AddLog(log string) {
	ic.Ic(s.Name, log)
	s.Logs = [10]string{log, s.Logs[0], s.Logs[1], s.Logs[2], s.Logs[3], s.Logs[4], s.Logs[5], s.Logs[6], s.Logs[7], s.Logs[8]}
}

func (s *Bot) GetLogs() string {
	logs := ""
	for _, log := range s.Logs {
		logs += log + "\n"
	}
	bts, _ := json.Marshal(strings.TrimSpace(logs))
	return string(bts)
}

func handleConnection(conn net.Conn) {
	var bot *Bot
	/* trunk-ignore(golangci-lint/ineffassign) */

	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		ic.Ic(err)
		debug.PrintStack()
		return
	}
	splits := bytes.Split(data, []byte("\x00"))
	botname := string(splits[0])
	//check if botname already exists
	ic.Ic(botname)
	if mbot, ok := Bots[botname]; ok {
		bot = mbot
		ic.Ic("already exists")
	} else {
		bot = &Bot{
			Name:      botname,
			Logs:      [10]string{},
			Timestamp: time.Now().Unix(),
		}
		Bots[bot.Name] = bot
	}
	for i := 1; i < len(splits); i++ {
		if len(splits[i]) > 0 {
			bot.AddLog(string(splits[i]))

		}

	}
	bot.Timestamp = time.Now().Unix()

	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			//check err type
			if err.Error() == "EOF" {
				ic.Ic("EOF", bot.Name)
			} else {
				ic.Ic(err)
				debug.PrintStack()
			}
			return
		}
		splits := bytes.Split(data, []byte("\x00"))
		for i := 0; i < len(splits); i++ {
			if len(splits[i]) > 0 {

				bot.AddLog(string(splits[i]))

			}
		}
		bot.Timestamp = time.Now().Unix()
	}
}

func TcpLogsServer() {

	ln, err := net.Listen("tcp", ":12225")
	if err != nil {
		ic.Ic(err)
		debug.PrintStack()
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			ic.Ic(err)
			debug.PrintStack()
			continue
		}
		go handleConnection(conn)
	}
}

func RunShell(name string, arg ...string) int {
	c := exec.Command(name, arg...)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		return 1
	}
	return 0
}

type apiHandler struct{}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/getLogs":
		data, _ := json.Marshal(Bots)
		//remove last char
		fmt.Fprint(w, string(data))
	case "/reset":
		ic.Ic("reset")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		ic.Ic(name)
		RunShell("bash", "scripts/reset.sh", name)

	default:
		fmt.Fprint(w, "hello world")

	}
}
func HttpLogsServer() {

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/api/", http.StripPrefix("/api", &apiHandler{}))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	ic.Ic(srv.ListenAndServe())
}
func main() {
	args := os.Args
	ic.Ic(args)
	switch len(args) {
	case 1:
		go TcpLogsServer()
		HttpLogsServer()
	}
}

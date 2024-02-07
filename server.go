package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strconv"
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
type FilderData struct {
	Name      string
	Text      string
	BotsCount int
}

func log(text string) {
	fmt.Println(text)
	conn, err := net.Dial("tcp", "127.0.0.1:12225")
	if err != nil {
		fmt.Println("Ошибка при подключении:", err)
		return
	}
	defer conn.Close()

	// Отправляем "server_log" и "INFO" раздельно с нулевым байтом в конце каждого сообщения
	conn.Write([]byte("server_log\000"))
	conn.Write([]byte("INFO " + text + "\000"))
}

func loge(text string) {
	fmt.Println(text)
	conn, err := net.Dial("tcp", "127.0.0.1:12225")
	if err != nil {
		fmt.Println("Ошибка при подключении:", err)
		return
	}
	defer conn.Close()

	// Отправляем "server_log" и "INFO" раздельно с нулевым байтом в конце каждого сообщения
	conn.Write([]byte("server_log\000"))
	conn.Write([]byte("ERROR " + text + "\000"))
}
func (s *Bot) Json() string {
	bts, _ := json.Marshal(s)
	return string(bts)

}
func (s *Bot) AddLog(log string) {
	ic.Ic(s.Name, log)
	if s.Name == "server_log" {
		s.Logs[0] = s.Logs[0] + "\n" + log
		return
	}
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

func GetOut(name string, arg ...string) string {
	c := exec.Command(name, arg...)
	stderr, _ := c.StderrPipe()
	stdout, _ := c.StdoutPipe()
	scannere := bufio.NewScanner(stderr)
	scannero := bufio.NewScanner(stdout)
	res := ""
	if err := c.Start(); err != nil {
		fmt.Println("Error: ", err)
	}
	for scannere.Scan() {
		res += scannere.Text()
	}
	for scannero.Scan() {
		res += scannero.Text()
	}
	return res

}
func RunShell(name string, arg ...string) int {
	log("Running: " + name + " " + strings.Join(arg, " "))
	c := exec.Command(name, arg...)
	stderr, _ := c.StderrPipe()
	stdout, _ := c.StdoutPipe()
	scannere := bufio.NewScanner(stderr)
	scannero := bufio.NewScanner(stdout)
	if err := c.Start(); err != nil {
		fmt.Println("Error: ", err)
		ic.Ic(err)
		for scannere.Scan() {
			loge(scannere.Text())
		}
		for scannero.Scan() {
			log(scannero.Text())
		}
		return 1
	}
	for scannere.Scan() {
		log(scannere.Text())
	}
	for scannero.Scan() {
		log(scannero.Text())
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
	case "/sigr":
		ic.Ic("sigr")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		ic.Ic(name)
		RestartBot(name)
		fmt.Fprint(w, "OK")

	case "/sigd":
		ic.Ic("sigd")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		ic.Ic(name)
		RunShell("bash", "scripts/csig.sh", name)
	case "/siga":
		ic.Ic("siga")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		ic.Ic(body)
		if name == "all" {
			RestartBots()
			fmt.Fprint(w, "OK")
		} else {
			ic.Ic(name)
			fmt.Fprint(w, "OK")
			RunShell("bash", "scripts/asig.sh", name)
		}
	case "/autoReg":
		ic.Ic("autoReg")
		body, _ := io.ReadAll(req.Body)
		data := strings.Split(string(body), "\n")
		prefix, count, folder, tables := data[0], data[1], data[2], data[3]
		ic.Ic(prefix, count, folder, tables)
		RunShell("/usr/local/bin/python", "scripts/autoreg.py", prefix, count, folder, tables)
		fmt.Fprint(w, "OK")

	case "/getFolders":
		Fdatas := []FilderData{}
		ic.Ic("getFolders")
		files, _ := os.ReadDir("../data")
		for _, file := range files {
			tmp := FilderData{}
			ic.Ic(file.Name())
			tmp.Name = file.Name()
			tmp.BotsCount = 0
			files2, _ := os.ReadDir("../data/" + file.Name())
			for _, file2 := range files2 {
				if file2.Name() == "text.cfg" {
					text, _ := os.ReadFile("../data/" + file.Name() + "/" + file2.Name())
					tmp.Text = string(text)

				} else if strings.HasPrefix(file2.Name(), "start_") && strings.HasSuffix(file2.Name(), ".sh") {
					tmp.BotsCount += 1
				}
			}
			Fdatas = append(Fdatas, tmp)
		}
		data, _ := json.Marshal(Fdatas)
		fmt.Fprint(w, string(data))

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

func FindBotContainer(name string) (string, int) {
	docOut := GetOut("docker", "ps", "-a", "--format", "&{{.Names}}")

	for _, namep := range strings.Split(docOut, "&") {

		docOut := GetOut("docker", "exec", "-i", namep, "ps", "--no-headers", "-eo", "pid,cmd")
		re := regexp.MustCompile(`(?P<first>\d+) python main.py ` + name + " " + name)
		_id := re.FindStringSubmatch(docOut)
		if len(_id) > 0 {
			id, _ := strconv.Atoi(_id[1])
			ic.Ic(id)
			return namep, id
		}

	}
	return "", 0
}

func KillBot(pid int, cont string) {
	RunShell("docker", "exec", cont, "kill", "-9", strconv.Itoa(pid))

}
func StartBot(name, cont string) {
	RunShell("docker", "exec", "-d", cont, "bash", "start_"+name+".sh")
}
func RestartBot(name string) {
	cont, id := FindBotContainer(name)
	ic.Ic(cont, id)
	KillBot(id, cont)
	StartBot(name, cont)
}
func RestartBots() {
	docOut := GetOut("docker", "ps", "-a", "--format", "&{{.ID}}*{{.Image}}")
	lines := strings.Split(docOut, "&")
	for _, line := range lines {
		if line == "" {
			continue
		}
		ic.Ic(line)
		fields := strings.Split(line, "*")
		id, image := fields[0], fields[1]
		ic.Ic(id, image)
		if image == "pp" {
			RunShell("docker", "restart", id)
		}
	}

}
func main() {
	args := os.Args
	ic.Ic(args)
	switch len(args) {
	case 1:
		go TcpLogsServer()
		time.Sleep(1 * time.Second)
		go log("server started")
		HttpLogsServer()
	default:
		RestartBot("sssspchz1xg")
	}
}

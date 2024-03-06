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
	fmt.Println(s.Name, log)
	// if s.Name == "server_log" {
	// 	s.Logs[0] = s.Logs[0] + "\n" + log
	// 	return
	// }
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
		fmt.Println(err)
		debug.PrintStack()
		return
	}
	splits := bytes.Split(data, []byte("\x00"))
	botname := string(splits[0])
	//check if botname already exists
	fmt.Println(botname)
	if mbot, ok := Bots[botname]; ok {
		bot = mbot
		fmt.Println("already exists")
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
				fmt.Println("EOF", bot.Name)
			} else {
				fmt.Println(err)
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
		fmt.Println(err)
		debug.PrintStack()
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
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

func CheckProxy(proxy_list string) [][]string {
	//uri = https://api.proxy-checker.net/api/proxy-checker/

	// filename := "prox.txt"
	// filedata, _ := os.ReadFile(filename)
	data := []byte(`proxy_list=` + strings.ReplaceAll(strings.ReplaceAll(proxy_list, "\n", "%0A"), ":", "%3A"))
	req, err := http.NewRequest("POST", "https://api.proxy-checker.net/api/proxy-checker/", bytes.NewReader(data))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.Status)
			body, _ := io.ReadAll(resp.Body)
			reg := regexp.MustCompile(`"initial": "([^,]+)", "valid": ([^,]+)`)
			res := reg.FindAllStringSubmatch(string(body), -1)
			return res
		}

	}
	return [][]string{}
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
		fmt.Println(err)
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
		fmt.Println("sigr")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		fmt.Println(name)
		RestartBot(name)
		fmt.Fprint(w, "OK")

	case "/sigd":
		fmt.Println("sigd")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		fmt.Println(name)
		DeleteBot(name)
		// RunShell("bash", "scripts/csig.sh", name)
	case "/siga":
		fmt.Println("siga")
		body, _ := io.ReadAll(req.Body)
		name := string(body)
		fmt.Println(body)
		if name == "all" {
			RestartBots()
			fmt.Fprint(w, "OK")
		} else {
			fmt.Println(name)
			fmt.Fprint(w, "OK")
			AutoReplase(name)
			//RunShell("bash", "scripts/asig.sh", name)
		}
	case "/sigclear":
		fmt.Println("sigclear")
		clearCache()
		log("cache cleared")
		fmt.Fprint(w, "OK")

	case "/autoReg":
		fmt.Println("autoReg")
		body, _ := io.ReadAll(req.Body)
		data := strings.Split(string(body), "\n")
		prefix, count, folder, tables := data[0], data[1], data[2], data[3]
		fmt.Println(prefix, count, folder, tables)
		RunShell("/usr/local/bin/python", "scripts/autoreg.py", prefix, count, folder, tables)
		fmt.Fprint(w, "OK")

	case "/getProxyes":
		fmt.Println("getProxyes")
		data, err := os.ReadFile("prox.txt")
		if err != nil {
			w.Write([]byte(err.Error()))
			//send 500 status code
			w.WriteHeader(500)
		} else {
			w.Write(data)
		}
	case "/checkProxy":
		fmt.Println("checkProxy")
		data, err := io.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			res := CheckProxy(string(data))
			for _, sr := range res {
				w.Write([]byte(sr[2]))
				w.Write([]byte(" "))
				w.Write([]byte(sr[1]))
				w.Write([]byte("\n"))
			}

		}
	case "/updateProxy":
		fmt.Println("updateProxy")
		data, err := io.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		os.WriteFile("prox.txt", data, 0644)

		w.Write([]byte(GetOut("/usr/local/bin/python", "scripts/updateProx.py")))

	case "/getFolders":
		Fdatas := []FilderData{}
		fmt.Println("getFolders")
		files, _ := os.ReadDir("../data")
		for _, file := range files {
			tmp := FilderData{}
			fmt.Println(file.Name())
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
	fmt.Println(srv.ListenAndServe())
}
func AutoReplase(name string) {
	data, folder := DeleteBot(name)
	re := regexp.MustCompile(`[^ ]+ [^ ]+ [^ ]+ [^ ]+ "([^"]+)"`)
	tables := re.FindStringSubmatch(data)[1]
	currentTime := time.Now()
	out := GetOut("/usr/local/bin/python", "scripts/autoreg.py", "replaсe"+currentTime.Format("2006-01-02"), "1", folder, tables)
	re = regexp.MustCompile(`{([^}]+)}`)
	logins := re.FindStringSubmatch(out)
	for _, cont := range GetContainers() {
		if strings.HasPrefix(cont, "pp") {
			StartBot(logins[1], cont)
		}
	}
}

func DeleteBot(name string) (string, string) {
	files, _ := os.ReadDir("/data")
	cont, id := FindBotContainer(name)
	KillBot(id, cont)
	for bot := range Bots {
		if bot == name {
			delete(Bots, name)
		}
	}
	for _, file := range files {
		if file.IsDir() {
			ffiles, _ := os.ReadDir("/data/" + file.Name())
			for _, ffile := range ffiles {
				if ffile.Name() == "start_"+name+".sh" {
					//remove file
					file_data, _ := os.ReadFile("/data/" + file.Name() + "/" + ffile.Name())
					os.Remove("/data/" + file.Name() + "/" + ffile.Name())
					return string(file_data), file.Name()
				}
			}
		}
	}
	return "", ""
}
func FindBotContainer(name string) (string, int) {

	for _, namep := range GetContainers() {

		docOut := GetOut("docker", "exec", "-i", namep, "ps", "--no-headers", "-eo", "pid,cmd")
		re := regexp.MustCompile(`(?P<first>\d+) python main.py ` + name + " " + name)
		_id := re.FindStringSubmatch(docOut)
		if len(_id) > 0 {
			id, _ := strconv.Atoi(_id[1])
			fmt.Println(id)
			return namep, id
		}

	}
	return "", 0
}

func GetContainers() []string {
	docOut := GetOut("docker", "ps", "-a", "--format", "&{{.Names}}")
	return strings.Split(docOut, "&")
}

func KillBot(pid int, cont string) {
	RunShell("docker", "exec", cont, "kill", "-9", strconv.Itoa(pid))

}
func StartBot(name, cont string) {
	RunShell("docker", "exec", "-d", cont, "bash", "start_"+name+".sh")
}
func RestartBot(name string) {
	cont, id := FindBotContainer(name)
	fmt.Println(cont, id)
	KillBot(id, cont)
	StartBot(name, cont)
}

func clearCache() {
	files, _ := os.ReadDir("/data")
	for _, file := range files {
		if file.IsDir() {
			ffiles, _ := os.ReadDir("/data/" + file.Name())
			for _, ffile := range ffiles {
				if ffile.Name() == "cache" {
					//remove all files in cache folder
					cacheFiles, _ := os.ReadDir("/data/" + file.Name() + "/cache")
					for _, cacheFile := range cacheFiles {
						os.Remove("/data/" + file.Name() + "/cache/" + cacheFile.Name())
					}
				}
			}
		}

	}
}
func RestartBots() {
	docOut := GetOut("docker", "ps", "-a", "--format", "&{{.ID}}*{{.Image}}")
	lines := strings.Split(docOut, "&")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fmt.Println(line)
		fields := strings.Split(line, "*")
		id, image := fields[0], fields[1]
		fmt.Println(id, image)
		if strings.HasPrefix(image, "pp") {
			RunShell("docker", "restart", id)
		}
	}

}
func main() {
	args := os.Args
	fmt.Println(args)
	switch len(args) {
	case 1:
		go TcpLogsServer()
		time.Sleep(1 * time.Second)
		go log("server started")
		HttpLogsServer()

	default:
		data := `python main.py "testdjf2l60gi2" "testdjf2l60gi2" "50K 500K 1M 5M 20M" "80.243.133.41:8000:U5apE5:9aETez" "45.132.21.144:8000:1znyDA:AfmfRG"
		`
		re := regexp.MustCompile(`[^ ]+ [^ ]+ [^ ]+ [^ ]+ "([^"]+)"`)
		tables := re.FindStringSubmatch(data)
		fmt.Println(tables)
		fmt.Println(tables[0])
		fmt.Println(tables[1])
		folder := "bot6"
		a := GetOut("/usr/local/bin/python", "scripts/autoreg.py", "replase", "1", folder, tables[1])
		re = regexp.MustCompile(`{([^}]+)}`)
		logins := re.FindStringSubmatch(a)

		fmt.Println(logins[1])
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	ServerRoot string
	Debug      bool
}

// InitServer tries to open "luckynumber.conf" in current directory
// and read it's properties. Properties are separated by whitespace.
// Currently the following properties are valid:
//   ServerRoot - specifies the root path for serving files
//   Debug      - turn debug output on or off
func InitServer() *Config {
	file, err := os.Open("luckynumber.conf")
	if err != nil {
		return &Config{}
	}
	defer file.Close()

	config := &Config{}

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		} // ignore lines starting with '#'
		tokens := make([]string, 0)
		tokens = strings.Split(line, " ")
		if tokens[0] == "ServerRoot" {
			// check if value starts and ends with "
			if len(tokens[1]) > 2 && tokens[1][0] == '"' && tokens[1][len(tokens[1])-1] == '"' {
				config.ServerRoot = tokens[1][1 : len(tokens[1])-1]
				if strings.HasSuffix(config.ServerRoot, "/") == false {
					config.ServerRoot = config.ServerRoot + "/"
				}
			}
			continue
		}
		if tokens[0] == "Debug" {
			if tokens[1] == "1" || tokens[1] == "True" || tokens[1] == "TRUE" {
				config.Debug = true
			}
			continue
		}

	}
	return config
}

type WebHandler struct {
	config Config
}

func (hdl WebHandler) HandleAjaxRequest(w http.ResponseWriter) {
	time.Sleep(3 * time.Second)
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(42) + 1 // 1 ... 42
	fmt.Fprintf(w, "Your lucky number is: %d", num)
}

func (hdl WebHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if hdl.config.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	}

	path := req.URL.Path[1:]
	if path == "" {
		path = "start.html"
	}
	if path == "getnumber" {
		hdl.HandleAjaxRequest(w)
		return
	}

	file, err := os.Open(hdl.config.ServerRoot + path)
	if err != nil {
		w.WriteHeader(404)
		io.WriteString(w, "Page not found")
		return
	}
	defer file.Close()
	rdr := bufio.NewReader(file)
	mimetype := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Add("Content-Type", mimetype)
	rdr.WriteTo(w)
}

func main() {
	config := InitServer()
	MyHandler := WebHandler{config: *config}
	fmt.Println("Server root: ", config.ServerRoot)
	fmt.Println("Debug mode: ", config.Debug)
	http.Handle("/", MyHandler)
	http.ListenAndServe(":8081", nil)

}

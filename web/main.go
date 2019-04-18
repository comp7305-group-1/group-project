package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type StandardOutputStreams struct {
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

type StringChannels struct {
	ChStdout <-chan string
	ChStderr <-chan string
}

var (
	processMapLock sync.Mutex
	processMap     = map[string]*StandardOutputStreams{}
	processMapAlt  = map[string]StringChannels{}
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/SparkPi", handleSparkPi)
	http.HandleFunc("/SparkPiSubmit", handleSparkPiSubmit)
	http.HandleFunc("/SparkPiSubmitStdoutWS", handleSparkPiSubmitStdoutWS)
	http.HandleFunc("/SparkPiSubmitStderrWS", handleSparkPiSubmitStderrWS)
	http.HandleFunc("/SparkPiSubmitResult", handleSparkPiSubmitResult)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func writeError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte(err.Error()))
}

func writeErrorString(w http.ResponseWriter, s string) {
	log.Println(s)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte(s))
}

func writeFile(w http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		writeError(w, err)
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Write(bytes)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "index.html")
}

func handleSparkPi(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "spark_pi.html")
}

func handleSparkPiSubmit(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("spark_pi_submit.gohtml")
	if err != nil {
		writeError(w, err)
		return
	}

	var requestID string
	var cmd *exec.Cmd

	result := func() bool {
		processMapLock.Lock()
		defer processMapLock.Unlock()
		for {
			requestIDInt := rand.Int()
			requestID = strconv.Itoa(requestIDInt)
			fmt.Printf("requestID: %s\n", requestID)
			if _, ok := processMap[requestID]; !ok {
				break
			}
		}

		cmd = exec.Command("ls", "-al", "/sbin")

		// stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			writeError(w, err)
			return false
		}

		// stderr
		stderr, err := cmd.StderrPipe()
		if err != nil {
			writeError(w, err)
			return false
		}

		processMap[requestID] = &StandardOutputStreams{stdout, stderr}

		return true
	}()

	if !result {
		return
	}

	if err := cmd.Start(); err != nil {
		writeError(w, err)
		return
	}
	fmt.Println("cmd.Start() returned")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, struct{ RequestID string }{requestID})
	fmt.Println("t.Execute() returned")
	w.(http.Flusher).Flush()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("cmd.Wait() returned")
}

func handleSparkPiSubmitStdoutWS(w http.ResponseWriter, req *http.Request) {
	url := req.URL
	if url == nil {
		writeErrorString(w, "url is nil")
		return
	}
	query := url.Query()
	if query == nil {
		writeErrorString(w, "query is nil")
		return
	}
	if rids, ok := query["rid"]; !ok {
		writeErrorString(w, "rid is not a key of query")
		return
	} else if len(rids) == 0 {
		writeErrorString(w, "len(rid) == 0")
		return
	} else if rid := rids[0]; len(rid) == 0 {
		writeErrorString(w, "len(rid[0]) == 0")
		return
	} else if streams, ok := processMap[rid]; !ok {
		writeErrorString(w, "invalid rid: "+rid)
		return
	} else {
		fmt.Printf("rid = %s\n", rid)
		conn, err := websocket.Upgrade(w, req, w.Header(), 1024, 1024)
		if err != nil {
			writeError(w, err)
			return
		}

		go func() {
			defer conn.Close()

			stdout := streams.Stdout
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				// if err := conn.WriteMessage(websocket.TextMessage, buffer); err != nil {
				// 	fmt.Println(err)
				// 	break
				// }
			}
		}()
	}
}

func handleSparkPiSubmitStderrWS(w http.ResponseWriter, req *http.Request) {
}

func handleSparkPiSubmitResult(w http.ResponseWriter, req *http.Request) {
	// TODO: Write this!
}

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

func main0() {
	cmd := exec.Command("ls", "-al", "/sbin")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("STDERR: " + line)
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

var (
	processMapLock sync.Mutex
	processMap     = map[string]struct{ Stdout, Stderr io.ReadCloser }{}
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/SparkPi", handleSparkPi)
	http.HandleFunc("/SparkPiSubmit", handleSparkPiSubmit)
	http.HandleFunc("/SparkPiSubmitStdout", handleSparkPiSubmitStdout)
	http.HandleFunc("/SparkPiSubmitStderr", handleSparkPiSubmitStderr)
	http.HandleFunc("/SparkPiSubmitResult", handleSparkPiSubmitResult)
	http.HandleFunc("/SlowResponseTest", handleSlowResponseTest)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte(err.Error()))
}

func writeErrorString(w http.ResponseWriter, s string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte(s))
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("index.html")
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

func handleSparkPi(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("spark_pi.html")
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

func handleSparkPiSubmit(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("spark_pi_submit.gohtml")
	if err != nil {
		writeError(w, err)
		return
	}
	requestID := 0
	requestIDString := ""
	func() {
		processMapLock.Lock()
		defer processMapLock.Unlock()
		for {
			requestID = rand.Int()
			requestIDString = strconv.Itoa(requestID)
			fmt.Println(requestIDString)
			if _, ok := processMap[requestIDString]; !ok {
				break
			}
			break
		}
		processMap[requestIDString] = struct{ Stdout, Stderr io.ReadCloser }{}
	}()

	cmd := exec.Command("ls", "-al", "/sbin")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	processMap[requestIDString] = struct{ Stdout, Stderr io.ReadCloser }{stdout, stderr}

	if err := cmd.Start(); err != nil {
		writeError(w, err)
		return
	}
	// fmt.Println("cmd.Start() called")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, struct{ RequestID string }{requestIDString})
	// fmt.Println("t.Execute() called")

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func handleSparkPiSubmitStdout(w http.ResponseWriter, req *http.Request) {
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
	} else if readers, ok := processMap[rid]; !ok {
		writeErrorString(w, "invalid rid: "+rid)
		return
	} else {
		// upgrader := websocket.Upgrader{}
		conn, err := websocket.Upgrade(w, req, w.Header(), 4096, 4096)
		if err != nil {
			// writeError(w, err)
			log.Println(err)
			return
		}
		defer conn.Close()

		writeCloser, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			// writeError(w, err)
			log.Println(err)
			return
		}

		stdout := readers.Stdout
		// w.WriteHeader(http.StatusOK)
		// w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		// io.Copy(writeCloser, stdout)
		buffer := make([]byte, 4096)
		for _, err := stdout.Read(buffer); err != nil; {
			if _, err2 := writeCloser.Write(buffer); err2 != nil {
				break
			}
		}
	}
}

func handleSparkPiSubmitStderr(w http.ResponseWriter, req *http.Request) {
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
	} else if readers, ok := processMap[rid]; !ok {
		writeErrorString(w, "invalid rid: "+rid)
		return
	} else {
		stderr := readers.Stderr
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		io.Copy(w, stderr)
	}
}

func handleSparkPiSubmitResult(w http.ResponseWriter, req *http.Request) {
	// TODO: Write this!
}

func handleSlowResponseTest(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	if f, ok := w.(http.Flusher); ok {
		fmt.Println("Yes")
		fmt.Fprintf(w, "%d\n", rand.Int())
		f.Flush()
		for i := 0; i < 100; i++ {
			fmt.Fprintf(w, "%d\n", i)
			f.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	} else {
		fmt.Fprintf(w, "%d\n", rand.Int())
		for i := 0; i < 100; i++ {
			fmt.Fprintf(w, "%d\n", i)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

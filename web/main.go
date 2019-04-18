package main

import (
	"bufio"
	"fmt"
	"html/template"
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

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/SparkPi", handleSparkPi)
	http.HandleFunc("/SparkPiSubmit", handleSparkPiSubmit)
	http.HandleFunc("/SparkPiSubmitWS", handleSparkPiSubmitWS)
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

	requestIDInt := rand.Int()
	requestID := strconv.Itoa(requestIDInt)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, struct{ RequestID string }{requestID})
}

func handleSparkPiSubmitWS(w http.ResponseWriter, req *http.Request) {
	// url := req.URL
	// if url == nil {
	// 	writeErrorString(w, "url is nil")
	// 	return
	// }
	// query := url.Query()
	// if query == nil {
	// 	writeErrorString(w, "query is nil")
	// 	return
	// }
	// if rids, ok := query["rid"]; !ok {
	// 	writeErrorString(w, "rid is not a key of query")
	// 	return
	// } else if len(rids) == 0 {
	// 	writeErrorString(w, "len(rid) == 0")
	// 	return
	// } else if rid := rids[0]; len(rid) == 0 {
	// 	writeErrorString(w, "len(rid[0]) == 0")
	// 	return
	// } else if streams, ok := processMap[rid]; !ok {
	// 	writeErrorString(w, "invalid rid: "+rid)
	// 	return
	// }

	// cmd
	cmd := exec.Command("ls", "-alR", "/etc/")

	// stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	// stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	conn, err := websocket.Upgrade(w, req, w.Header(), 1024, 1024)
	if err != nil {
		writeError(w, err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			if err := conn.WriteJSON(map[string]string{"stream": "stdout", "message": line}); err != nil {
				fmt.Println(err)
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			if err := conn.WriteJSON(map[string]string{"stream": "stderr", "message": line}); err != nil {
				fmt.Println(err)
				break
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		writeError(w, err)
		return
	}
	// fmt.Println("cmd.Start() returned")

	wg.Wait()
	// fmt.Println("wg.Wait() returned")

	cmd.Wait()
	// fmt.Println("cmd.Wait() returned")
}

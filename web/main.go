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

	if err := req.ParseForm(); err != nil {
		writeError(w, err)
		return
	}

	partition := req.FormValue("pc")
	if _, err := strconv.Atoi(partition); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, struct{ Partition string }{partition})
}

func handleSparkPiSubmitWS(w http.ResponseWriter, req *http.Request) {
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
	partitionSlice, ok := query["pc"]
	if !ok {
		writeErrorString(w, "pc is not a key of query")
		return
	}
	if len(partitionSlice) == 0 {
		writeErrorString(w, "len(partitionSlice) == 0")
		return
	}
	partition := partitionSlice[0]
	if _, err := strconv.Atoi(partition); err != nil {
		writeError(w, err)
		return
	}

	// Command
	// cmd := exec.Command("ls", "-alR", "/etc/")
	cmd := exec.Command("spark-submit", "--class", "org.apache.spark.examples.SparkPi", "--master", "yarn", "--deploy-mode", "client", "/opt/spark-2.4.0-bin-hadoop2.7/examples/jars/spark-examples_2.11-2.4.0.jar", partition)

	// Standard Output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	// Standard Error
	stderr, err := cmd.StderrPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	// WebSocket Connection
	conn, err := websocket.Upgrade(w, req, w.Header(), 1024, 1024)
	if err != nil {
		writeError(w, err)
		return
	}
	defer conn.Close()

	// WaitGroup
	var wg sync.WaitGroup

	// Goroutine for handling stdout
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

	// Goroutine for handling stderr
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

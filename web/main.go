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
	http.HandleFunc("/SparkPiResult", handleSparkPiResult)
	http.HandleFunc("/SparkPiResultWS", handleSparkPiResultWS)

	http.HandleFunc("/Mysteries", handleMysteries)
	http.HandleFunc("/MysteriesResult", handleMysteriesResult)
	http.HandleFunc("/MysteriesResultWS", handleMysteriesResultWS)

	http.HandleFunc("/MysteriesStream", handleMysteriesStream)

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

func handleSparkPiResult(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("spark_pi_result.gohtml")
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

func handleSparkPiResultWS(w http.ResponseWriter, req *http.Request) {
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

func handleMysteries(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "mysteries.html")
}

func handleMysteriesResult(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("mysteries_result.gohtml")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		writeError(w, err)
		return
	}

	mysteryText := req.FormValue("mt")

	mode := req.FormValue("mode")
	var isUsedCachedResult string
	switch mode {
	case "1":
		isUsedCachedResult = "No"
		break
	case "2":
		isUsedCachedResult = "Yes"
		break
	default:
		writeErrorString(w, "invalid mode")
		return
	}

	partitionCount := req.FormValue("pc")
	if _, err := strconv.Atoi(partitionCount); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, map[string]string{
		"MysteryText":        mysteryText,
		"Mode":               mode,
		"IsUsedCachedResult": isUsedCachedResult,
		"PartitionCount":     partitionCount,
	})
}

func handleMysteriesResultWS(w http.ResponseWriter, req *http.Request) {
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

	// Mystery Text
	mysteryTextSlice, ok := query["mt"]
	if !ok {
		writeErrorString(w, "mt is not a key of query")
		return
	}
	if len(mysteryTextSlice) == 0 {
		writeErrorString(w, "len(mysteryTextSlice) == 0")
		return
	}
	mysteryText := mysteryTextSlice[0]
	if len(mysteryText) == 0 {
		writeErrorString(w, "len(mysteryText) == 0")
		return
	}

	// Mode
	modeSlice, ok := query["mode"]
	if !ok {
		writeErrorString(w, "mode is not a key of query")
		return
	}
	if len(modeSlice) == 0 {
		writeErrorString(w, "len(modeSlice) == 0")
		return
	}
	mode := modeSlice[0]
	if mode != "0" && mode != "1" && mode != "2" {
		writeErrorString(w, "invalid mode")
		return
	}

	// Partition Count
	partitionCountSlice, ok := query["pc"]
	if !ok {
		writeErrorString(w, "pc is not a key of query")
		return
	}
	if len(partitionCountSlice) == 0 {
		writeErrorString(w, "len(partitionCountSlice) == 0")
		return
	}
	partitionCount := partitionCountSlice[0]
	if _, err := strconv.Atoi(partitionCount); err != nil {
		writeError(w, err)
		return
	}

	// Command
	var cmd *exec.Cmd
	switch mode {
	case "0", "1":
		cmd = exec.Command("spark-submit", "--master", "yarn", "--py-files", "../app/dependencies.zip", "../app/m4.py", mode, "hdfs://gpu1:8020/books", mysteryText, partitionCount, "1", "0")
	case "2":
		cmd = exec.Command("spark-submit", "--master", "yarn", "--py-files", "../app/dependencies.zip", "../app/m4.py", mode, mysteryText)
	}

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

func handleMysteriesStream(w http.ResponseWriter, req *http.Request) {
}

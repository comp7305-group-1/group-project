package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Main

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

// Utility Functions

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

// Index Page

func handleIndex(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "index.html")
}

// SparkPi

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

// Unravel the Mysteries by Spark

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

	// Mystery Text
	mysteryText := req.FormValue("mt")

	// Mode
	mode := req.FormValue("mode")
	var isUsedCachedResult string
	switch mode {
	case "0":
		isUsedCachedResult = "No (do not save the result to cache)"
	case "1":
		isUsedCachedResult = "No (save the result to cache)"
	case "2":
		isUsedCachedResult = "Yes (load the result from cache)"
	default:
		writeErrorString(w, "invalid mode")
		return
	}

	// Executor Count
	numExecutors := req.FormValue("ne")
	if _, err := strconv.Atoi(numExecutors); err != nil {
		writeError(w, err)
		return
	}

	// Executor Cores
	executorCores := req.FormValue("ec")
	if _, err := strconv.Atoi(executorCores); err != nil {
		writeError(w, err)
		return
	}

	// Driver Memory
	driverMemory := req.FormValue("dm")
	dmIndex := len(driverMemory) - 1
	driverMemoryNum := driverMemory[:dmIndex]
	driverMemoryUnit := driverMemory[dmIndex:]
	if _, err := strconv.Atoi(driverMemoryNum); err != nil {
		writeError(w, err)
		return
	}
	if driverMemoryUnit != "m" && driverMemoryUnit != "g" {
		writeErrorString(w, "the unit of dm must be 'm' or 'g'")
		return
	}

	// Partition Count
	minPartitions := req.FormValue("mp")
	if _, err := strconv.Atoi(minPartitions); err != nil {
		writeError(w, err)
		return
	}

	// Books Path
	booksPath := req.FormValue("bp")
	var booksPathURI string
	switch booksPath {
	case "0":
		booksPathURI = "har:///har/booksarchive.har"
	case "1":
		booksPathURI = "hdfs://gpu1:8020/books"
	default:
		writeErrorString(w, "invalid bp")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, map[string]string{
		"MysteryText":        mysteryText,
		"Mode":               mode,
		"IsUsedCachedResult": isUsedCachedResult,
		"NumExecutors":       numExecutors,
		"ExecutorCores":      executorCores,
		"DriverMemory":       driverMemory,
		"MinPartitions":      minPartitions,
		"BooksPath":          booksPath,
		"BooksPathURI":       booksPathURI,
	})
}

func getParam(w http.ResponseWriter, query url.Values, key string) (string, bool) {
	valueSlice, ok := query[key]
	if !ok {
		writeErrorString(w, fmt.Sprintf("%s is not a key of query", key))
		return "", false
	}
	if len(valueSlice) == 0 {
		writeErrorString(w, fmt.Sprintf("len(valueSlice) == 0, key = %s", key))
		return "", false
	}
	value := valueSlice[0]
	if len(value) == 0 {
		writeErrorString(w, fmt.Sprintf("len(value) == 0, key = %s", key))
		return "", false
	}
	return value, true
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
	mysteryText, ok := getParam(w, query, "mt")
	if !ok {
		return
	}

	// Mode
	mode, ok := getParam(w, query, "mode")
	if !ok {
		return
	}
	if mode != "0" && mode != "1" && mode != "2" {
		writeErrorString(w, "invalid mode")
		return
	}

	// Executor Count
	numExecutors, ok := getParam(w, query, "ne")
	if !ok {
		return
	}
	if _, err := strconv.Atoi(numExecutors); err != nil {
		writeError(w, err)
		return
	}

	// Executor Cores
	executorCores, ok := getParam(w, query, "ec")
	if !ok {
		return
	}
	if _, err := strconv.Atoi(executorCores); err != nil {
		writeError(w, err)
		return
	}

	// Driver Memory
	driverMemory, ok := getParam(w, query, "dm")
	if !ok {
		return
	}
	dmIndex := len(driverMemory) - 1
	driverMemoryNum := driverMemory[:dmIndex]
	driverMemoryUnit := driverMemory[dmIndex:]
	if _, err := strconv.Atoi(driverMemoryNum); err != nil {
		writeError(w, err)
		return
	}
	if driverMemoryUnit != "m" && driverMemoryUnit != "g" {
		writeErrorString(w, "the unit of dm must be 'm' or 'g'")
		return
	}

	// Partition Count
	minPartitions, ok := getParam(w, query, "mp")
	if !ok {
		return
	}
	if _, err := strconv.Atoi(minPartitions); err != nil {
		writeError(w, err)
		return
	}

	// Books Path
	booksPath, ok := getParam(w, query, "bp")
	if !ok {
		return
	}
	var booksPathURI string
	switch booksPath {
	case "0":
		booksPathURI = "har:///har/booksarchive.har"
	case "1":
		booksPathURI = "hdfs://gpu1:8020/books"
	default:
		writeErrorString(w, "invalid bp")
		return
	}

	// Command
	var cmd *exec.Cmd
	switch mode {
	case "0", "1":
		cmd = exec.Command("spark-submit", "--master", "yarn", "--num-executors", numExecutors, "--executor-cores", executorCores, "--driver-memory", driverMemory, "--py-files", "../app/dependencies.zip", "../app/m4.py", mode, booksPathURI, mysteryText, minPartitions, "1", "0")
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

// Unravel the Mysteries by Spark (Streaming Version)

func handleMysteriesStream(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "mysteries_stream.html")
}

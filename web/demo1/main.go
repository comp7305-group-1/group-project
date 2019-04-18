package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/WebSocketTest", handleWebSocketTest)
	http.HandleFunc("/WebSocketTestWS", handleWebSocketTestWS)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte(err.Error()))
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

func handleWebSocketTest(w http.ResponseWriter, req *http.Request) {
	writeFile(w, "websocket_test.html")
}

func handleWebSocketTestWS(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleWebSocketTestWS is called")
	conn, err := websocket.Upgrade(w, req, w.Header(), 1024, 1024)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go func() {
		defer conn.Close()
		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	fmt.Printf("Message Type: %d, Buffer: %s, Error: %s\n", messageType, string(p), err.Error())
		// 	return
		// } else {
		// 	fmt.Printf("Message Type: %d, Buffer: %s, Error: %s\n", messageType, string(p))
		// }
		fmt.Println("The for loop is about to start")
		for i := 0; i < 100; i++ {
			r := rand.Int()
			message := fmt.Sprintf("Hello, %d", r)
			// m := map[string]string{
			// 	"message": message,
			// }
			// conn.WriteJSON(m)
			conn.WriteMessage(websocket.TextMessage, []byte(message))
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("The for loop ended")
	}()
}

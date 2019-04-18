package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/SlowResponseTest", handleSlowResponseTest)
	log.Fatal(http.ListenAndServe(":12345", nil))
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
		fmt.Println("No")
		fmt.Fprintf(w, "%d\n", rand.Int())
		for i := 0; i < 100; i++ {
			fmt.Fprintf(w, "%d\n", i)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

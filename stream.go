package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func stream(w http.ResponseWriter, r *http.Request) {
	log.Println("Streaming")
	w.Header().Set("Content-Type", "text/plain")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for !done {
		fmt.Fprintf(w, "%d/%d\n", current_size, current_total_size)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("Streaming done!")
	fmt.Fprintf(w, "%d/%d\n", current_size, current_total_size)
	flusher.Flush()
}


package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StreamedData struct {
	CurrentSize      int64  `json:"currentSize"`
	CurrentTotalSize int64  `json:"currentTotalSize"`
	Videos           []Node `json:"videos"`
}

func stream(w http.ResponseWriter, r *http.Request) {
	log.Println("Streaming")
	w.Header().Set("Content-Type", "text/json")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	sd, _ := json.Marshal(StreamedData{
		current_size,
		current_total_size,
		videoQueue.GetAll(),
	})
	for !done {
		sd, err := json.Marshal(StreamedData{
			current_size,
			current_total_size,
			videoQueue.GetAll(),
		})

		if err != nil {
			log.Printf("failed to serialize json: %s: %s\n", sd, err.Error())
			break
		}

		// fmt.Fprintf(w, "%d/%d\n", current_size, current_total_size)
		fmt.Fprintf(w, "%s", sd)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("Streaming done!")
	fmt.Fprintf(w, "%s", sd)
	flusher.Flush()
}

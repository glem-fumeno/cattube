package main

import (
	"log"
	"net/http"
)

var (
	current_total_size int64
	current_size       int64
	done               bool
	videoQueue         Queue
)

func root(w http.ResponseWriter, r *http.Request) {
	log.Println("Accessing root")
	http.ServeFile(w, r, "./views/index.html")
}

func main() {
	videoQueue = MakeQueue()
	done = true
	log.Println("Starting server")
	http.HandleFunc("/", root)
	http.HandleFunc("POST /api/download", download)
	http.HandleFunc("GET /api/stream", stream)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":8090", nil)
}

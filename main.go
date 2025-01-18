package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	youtube "github.com/kkdai/youtube/v2"
)

var (
	current_total_size int64
	current_size       int64
	done               bool
)

func root(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/index.html")
}

type UrlRequest struct {
	Url string `json:"url"`
}

type Stream struct {
	stream io.Reader
}

func (s Stream) Read(p []byte) (n int, err error) {
	n, e := s.stream.Read(p)
	current_size += int64(n)
	return n, e
}

type Client struct {
	video  *youtube.Video
	client youtube.Client
}

func download_part(client Client, video_type string, part int) error {
	video_formats := client.video.Formats.Type(video_type)
	video_stream, video_size, err := client.client.GetStream(client.video, &video_formats[0])
	if err != nil {
		return err
	}
	defer video_stream.Close()
	video_part, err := os.Create(fmt.Sprintf("video.mp4.part%d", part))
	if err != nil {
		return err
	}
	defer video_part.Close()
	current_total_size = video_size
	current_size = 0
	_, err = io.Copy(video_part, Stream{video_stream})
	if err != nil {
		return err
	}
	return nil
}

func download(w http.ResponseWriter, r *http.Request) {
	done = false
	// current_total_size = 1
	// current_size = 0
	jsonParser := json.NewDecoder(r.Body)
	data := UrlRequest{}
	err := jsonParser.Decode(&data)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}
	client := youtube.Client{}

	video, err := client.GetVideo(data.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = download_part(Client{video, client}, "video/mp4", 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = download_part(Client{video, client}, "audio/mp4", 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "video.mp4.part1",
		"-i", "video.mp4.part2",
		"-c", "copy",
		"-shortest",
		"video.mp4",
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd = exec.Command(
		"mv",
		"video.mp4",
		fmt.Sprintf("%s [%s].mp4", video.Title, video.ID),
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message": "%s"}`, err)
		return
	}
	done = true
	fmt.Fprintf(w, `{"title": "%s"}`, video.Title)
}

func stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for !done {
		fmt.Fprintf(w, "%d/%d\n", current_size, current_total_size)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Fprintf(w, "%d/%d\n", current_size, current_total_size)
	flusher.Flush()
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("POST /api/download", download)
	http.HandleFunc("GET /api/stream", stream)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":8090", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	log.Println("Accessing root")
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
	log.Printf("Downloading part %d\n", part)
	video_formats := client.video.Formats.Type(video_type)
	video_stream, video_size, err := client.client.GetStream(client.video, &video_formats[0])
	if err != nil {
		return err
	}
	defer video_stream.Close()
	log.Printf("Got stream part %d\n", part)
	video_part, err := os.Create(fmt.Sprintf("/tmp/video.mp4.part%d", part))
	if err != nil {
		return err
	}
	defer video_part.Close()
	log.Printf("Created file part %d\n", part)
	current_total_size = video_size
	current_size = 0
	_, err = io.Copy(video_part, Stream{video_stream})
	log.Printf("Copied file part %d\n", part)
	if err != nil {
		return err
	}
	return nil
}

func merge_parts(title string, video_id string) error {
	log.Printf("Merging video and audio %s\n", title)
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "/tmp/video.mp4.part1",
		"-i", "/tmp/video.mp4.part2",
		"-c", "copy",
		"-shortest",
		"/tmp/video.mp4",
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Printf("Moving video %s\n", title)
	cmd = exec.Command(
		"mv",
		"/tmp/video.mp4",
		fmt.Sprintf("/destination/%s [%s].mp4", title, video_id),
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Printf("Cleaning up %s\n", title)
	cmd = exec.Command(
		"rm",
		"/tmp/video.mp4.part1",
		"/tmp/video.mp4.part2",
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func set_done() {
	done = true
}

func download(w http.ResponseWriter, r *http.Request) {
	done = false
	defer set_done()
	log.Printf("Decoding body")
	jsonParser := json.NewDecoder(r.Body)
	data := UrlRequest{}
	err := jsonParser.Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}
	client := youtube.Client{}
	log.Printf("Getting video %s\n", data.Url)

	video, err := client.GetVideo(data.Url)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = download_part(Client{video, client}, "video/mp4", 1)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = download_part(Client{video, client}, "audio/mp4", 2)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = merge_parts(video.Title, video.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Done downloading video %s\n", video.Title)
	fmt.Fprintf(w, `Done!`)
}

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

func main() {
	log.Println("Starting server")
	http.HandleFunc("/", root)
	http.HandleFunc("POST /api/download", download)
	http.HandleFunc("GET /api/stream", stream)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":8090", nil)
}

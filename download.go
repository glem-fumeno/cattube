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

func Log(message string, args ...interface{}) {
	current_log = fmt.Sprintf(message, args...)
	log.Println(current_log)
}

func download_part(client Client, video_type string, part int) error {
	Log("Downloading part %d", part)
	video_formats := client.video.Formats.Type(video_type)
	video_stream, video_size, err := client.client.GetStream(client.video, &video_formats[0])
	if err != nil {
		return err
	}
	defer video_stream.Close()
	Log("Got stream part %d", part)
	video_part, err := os.Create(fmt.Sprintf("/tmp/video.mp4.part%d", part))
	if err != nil {
		return err
	}
	defer video_part.Close()
	Log("Created file part %d", part)
	current_total_size = video_size
	current_size = 0
	_, err = io.Copy(video_part, Stream{video_stream})
	Log("Copied file part %d", part)
	if err != nil {
		return err
	}
	return nil
}

func merge_parts(title string, video_id string) (string, error) {
	Log("Merging video and audio %s", title)
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "/tmp/video.mp4.part1",
		"-i", "/tmp/video.mp4.part2",
		"-c", "copy",
		"-shortest",
		"/tmp/video.mp4",
	)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}
	Log("Moving video %s", title)
	cmd = exec.Command(
		"mv",
		"/tmp/video.mp4",
		fmt.Sprintf("/destination/%s [%s].mp4", title, video_id),
	)
	result, err = cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}
	Log("Cleaning up %s", title)
	cmd = exec.Command(
		"rm",
		"/tmp/video.mp4.part1",
		"/tmp/video.mp4.part2",
	)
	result, err = cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}
	return "", nil
}

func set_done() {
	if videoQueue.IsEmpty() {
		done = true
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	done = false
	defer set_done()
	Log("Decoding body")
	jsonParser := json.NewDecoder(r.Body)
	data := UrlRequest{}
	err := jsonParser.Decode(&data)
	if err != nil {
		Log(err.Error())
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}
	client := youtube.Client{}
	Log("Getting video %s", data.Url)

	video, err := client.GetVideo(data.Url)
	if err != nil {
		Log(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nodeData := &Node{data.Url, video.Title, video.Duration.String()}
	videoQueue.Enqueue(nodeData)
	defer videoQueue.Dequeue()
	val, _ := videoQueue.Peek()
	for val.Url != nodeData.Url {
		Log("Waiting for %v", val.Url)
		time.Sleep(1000 * time.Millisecond)
		val, _ = videoQueue.Peek()
	}
	err = download_part(Client{video, client}, "video/mp4", 1)
	if err != nil {
		Log(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = download_part(Client{video, client}, "audio/mp4", 2)
	if err != nil {
		Log(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := merge_parts(video.Title, video.ID)
	if err != nil {
		Log("Error merging parts: %s, %s", result, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Log("Done downloading video %s", video.Title)
	fmt.Fprintf(w, `Done!`)
}

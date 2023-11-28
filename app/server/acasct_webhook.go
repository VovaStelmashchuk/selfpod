package server

import (
	"encoding/json"
	"fmt"
	"io"
	"main/app/media"
	"net/http"
)

type Episode struct {
	Event       string `json:"event"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	PublishDate string `json:"publishDate"`
	CoverUrl    string `json:"coverUrl"`
	AudioUrl    string `json:"audioUrl"`
}

func AcastWebHook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["token"][0]

	fmt.Printf("Token: %s\n", query)

	var episode Episode

	err := json.NewDecoder(r.Body).Decode(&episode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Episode ID: %s, Title: %s\n", episode.ID, episode.Title)

	media.PrepareNewVideo(episode.AudioUrl, episode.CoverUrl)

	io.WriteString(w, "Received episode: "+episode.Title+"\n")
}

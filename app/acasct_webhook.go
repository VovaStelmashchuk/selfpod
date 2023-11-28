package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	hookToken := os.Getenv("HOOK_TOKEN")
	queryToken := r.URL.Query()["token"][0]

	if hookToken != queryToken {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	var episode Episode

	err := json.NewDecoder(r.Body).Decode(&episode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := CreateYoutubeItem(episode)

	fmt.Printf("%+v\n", res)

	io.WriteString(w, "Received episode: "+episode.ID+"\n")
}

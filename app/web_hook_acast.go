package app

import (
	"encoding/json"
	"io"
	"main/app/youtube_uploader"
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

	youtube_uploader.UploadToYoutube(
		youtube_uploader.YoutubeUploadRequset{
			Filename:    res.VideoFilePath,
			Title:       res.Title,
			Description: res.Description,
		},
	)

	io.WriteString(w, "Video uploaded to youtube")
}

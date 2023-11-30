package app

import (
	"encoding/json"
	"io"
	"main/app/config"
	"main/app/youtube_uploader"
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
	queryToken := r.URL.Query()["token"][0]

	if config.ACastHookToken != queryToken {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	episode, err := getEpisodeFromHook(w, r)

	res := CreateYoutubeItem(episode)

	youtube_uploader.UploadToYoutube(
		youtube_uploader.YoutubeUploadRequset{
			Filename:    res.VideoFilePath,
			Title:       res.Title,
			Description: res.Description,
		},
	)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	io.WriteString(w, "Video uploaded to youtube")
}

func getEpisodeFromHook(w http.ResponseWriter, r *http.Request) (Episode, error) {
	var episode Episode

	err := json.NewDecoder(r.Body).Decode(&episode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return episode, err
}

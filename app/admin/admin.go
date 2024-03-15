package admin

import (
	"encoding/json"
	"log"
	appconfig "main/app/config"
	"main/app/processing"
	"net/http"
	"strconv"
)

type EpisodeIdResponse struct {
	EpisodeId int `json:"episodeId"`
}

func TryUploadAgain(w http.ResponseWriter, r *http.Request) {
	queryToken := r.URL.Query()["token"][0]

	if appconfig.ACastHookToken != queryToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	episodeIdStr := r.URL.Query()["episodeId"][0]

	log.Printf("episode id %s", episodeIdStr)

	episodeId, err := strconv.Atoi(episodeIdStr)

	if err != nil {
		log.Printf("episode id not found")
		http.Error(w, "Episode id not provided", http.StatusBadRequest)
	}

	processing.ProcessEpisode(
		processing.ProcessEpisodeTask{
			EpisodeId: episodeId,
		},
	)

	response := EpisodeIdResponse{
		EpisodeId: episodeId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

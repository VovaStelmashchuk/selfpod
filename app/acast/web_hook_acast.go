package acast

import (
	"encoding/json"
	"log"
	"main/app/config"
	"main/app/database"
	"main/app/processing"
	"net/http"
)

type Episode struct {
	ID       string `json:"id"`
	CoverUrl string `json:"coverUrl"`
	AudioUrl string `json:"audioUrl"`
	Title    string `json:"title"`
}

type JsonResponse struct {
	EpisodeId int64 `json:"episode_id"`
}

func WebHook(w http.ResponseWriter, r *http.Request) {
	queryToken := r.URL.Query()["token"][0]

	if config.ACastHookToken != queryToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var episode Episode
	err := json.NewDecoder(r.Body).Decode(&episode)
	if err != nil {
		log.Printf("fail to decode body: %v", err)
		http.Error(w, "Fail to decode body", http.StatusInternalServerError)
	}

	id, err := database.SaveEpisode(
		database.Episode{
			AcastId:         episode.ID,
			Title:           episode.Title,
			ImageUrl:        episode.CoverUrl,
			AudioUrl:        episode.AudioUrl,
			ProcessingState: database.NOT_STARTED,
		},
	)

	processing.ProcessEpisode(
		processing.ProcessEpisodeTask{
			EpisodeId: int(id),
		},
	)

	response := JsonResponse{EpisodeId: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Printf("error when try to procces a cast web hoook: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

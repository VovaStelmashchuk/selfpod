package processing

import (
	"log"
	"main/app/database"
	"main/app/media"
)

type ProcessEpisodeTask struct {
	EpisodeId int
}

var taskQueue chan ProcessEpisodeTask

func init() {
	taskQueue = make(chan ProcessEpisodeTask, 10)

	go worker(taskQueue)
}

func ProcessEpisode(task ProcessEpisodeTask) {
	taskQueue <- task
}

func worker(taskQueue <-chan ProcessEpisodeTask) {
	for task := range taskQueue {
		processTask(task)
	}
}

func processTask(task ProcessEpisodeTask) {
	err := database.UpdateEpisodeState(task.EpisodeId, database.IN_PROGRESS)
	if err != nil {
		log.Printf("Error updating episode state to IN_PROGRESS: %v", err)
		return
	}

	episode, err := database.GetEpisode(task.EpisodeId)

	videoFile := media.PrepareNewVideo(episode.AudioUrl, episode.ImageUrl)

	episodeMetaInfo, err := GetEpisodeDescription(episode.AcastId)

	media.UploadToYoutube(
		media.YoutubeUploadRequset{
			Filename:    videoFile,
			Title:       episode.Title,
			Description: episodeMetaInfo + "\n Ви можете підтримати нас на https://www.patreon.com/androidstory",
		},
	)

	err = database.UpdateEpisodeState(task.EpisodeId, database.SUCCESS)

	if err != nil {
		log.Printf("Error updating episode state to SUCCESS: %v", err)
		database.UpdateEpisodeState(task.EpisodeId, database.FAIL)
	}
}

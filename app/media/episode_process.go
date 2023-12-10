package media

import (
	"log"
	"main/app/database"
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
	log.Printf("Processing episode %v", task.EpisodeId)
	err := database.UpdateEpisodeState(task.EpisodeId, database.IN_PROGRESS)
	if err != nil {
		log.Printf("Error updating episode state to IN_PROGRESS: %v", err)
		return
	}

	episode, err := database.GetEpisode(task.EpisodeId)

	videoFile := PrepareNewVideo(episode.AudioUrl, episode.ImageUrl)

	UploadToYoutube(
		YoutubeUploadRequset{
			Filename:    videoFile,
			Title:       episode.Title,
			Description: episode.Description + "\n Ви можете підтримати нас на https://www.patreon.com/androidstory",
		},
	)

	err = database.UpdateEpisodeState(task.EpisodeId, database.SUCCESS)

	if err != nil {
		log.Printf("Error updating episode state to SUCCESS: %v", err)
		database.UpdateEpisodeState(task.EpisodeId, database.FAIL)
	}
}

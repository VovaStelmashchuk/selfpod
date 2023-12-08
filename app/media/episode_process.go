package media

import (
	"fmt"
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
	fmt.Printf("Processing episode %v\n", task.EpisodeId)
	err := database.UpdateEpisodeState(task.EpisodeId, database.IN_PROGRESS)
	if err != nil {
		fmt.Print("Error updating episode state to IN_PROGRESS")
		return
	}

	episode, err := database.GetEpisode(task.EpisodeId)

	videoFile := PrepareNewVideo(episode.AudioUrl, episode.ImageUrl)

	UploadToYoutube(
		YoutubeUploadRequset{
			Filename:    videoFile,
			Title:       episode.Title,
			Description: episode.Description,
		},
	)

	if err != nil {
		database.UpdateEpisodeState(task.EpisodeId, database.FAIL)
	}
}

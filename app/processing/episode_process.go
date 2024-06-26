package processing

import (
	"log"
	"main/app/database"
	"main/app/media"
	"main/app/notifications"
	"time"
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
	log.Printf("Processing episode: %v", task.EpisodeId)
	err := database.UpdateEpisodeState(task.EpisodeId, database.InProgress)
	if err != nil {
		log.Printf("Error updating episode state to IN_PROGRESS: %v", err)
		return
	}

	episode, err := database.GetEpisode(task.EpisodeId)

	videoFile := media.PrepareNewVideo(episode.AudioUrl, episode.ImageUrl)

	log.Printf("Sleep for 10 seconds before getting episode description...")

	time.Sleep(10 * time.Second)

	episodeMetaInfo, err := GetEpisodeDescription(episode.AcastId)

	if err != nil {
		log.Printf("Error getting episode meta info: %v", err)
		_ = database.UpdateEpisodeState(task.EpisodeId, database.Fail)
		return
	}

	log.Printf("Episode meta info: %v", episodeMetaInfo)

	episodeDescription := episodeMetaInfo + "\n Ви можете підтримати нас на https://www.patreon.com/androidstory"

	err = database.AddDescriptionToEpisode(task.EpisodeId, episodeDescription)

	if err != nil {
		log.Printf("Error adding episode description: %v", err)
		return
	}

	err = media.UploadToYoutube(
		media.YoutubeUploadRequest{
			Filename:    videoFile,
			Title:       episode.Title,
			Description: episodeDescription,
		},
	)

	if err != nil {
		notifications.SendDiscordNotification("Error uploading episode to YouTube, check logs, episode title: " + episode.Title)
		return
	}

	err = database.UpdateEpisodeState(task.EpisodeId, database.Success)

	notifications.SendDiscordNotification("Episode uploaded: " + episode.Title)

	if err != nil {
		notifications.SendDiscordNotification("Something went wrong, check logs, episode title: " + episode.Title)
		log.Printf("Error updating episode state to Success: %v", err)
		_ = database.UpdateEpisodeState(task.EpisodeId, database.Fail)
	}
}

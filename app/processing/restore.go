package processing

import (
	"log"
	"main/app/database"
)

func RestoreFailedProcess() {
	log.Printf("Restoring failed processes...")
	episodeIds, err := database.GetEpisodeIdsByState(database.InProgress)

	if err != nil {
		log.Printf("Error while getting episodes: %v", err)
	}

	for _, episodeId := range episodeIds {
		log.Printf("Restoring episode %d", episodeId)
		ProcessEpisode(
			ProcessEpisodeTask{
				EpisodeId: episodeId,
			},
		)
	}
}

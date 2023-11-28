package app

type YoutubeItem struct {
	Title         string
	Description   string
	VideoFilePath string
}

func CreateYoutubeItem(episode Episode) YoutubeItem {
	//videoFile := media.PrepareNewVideo(episode.AudioUrl, episode.CoverUrl)
	videoFile := "tmp_files/output.mov"

	episodeMetaInfo := GetEpisodeMetaInfo(episode.ID)

	return YoutubeItem{
		Title:         episodeMetaInfo.Title,
		Description:   episodeMetaInfo.Description,
		VideoFilePath: videoFile,
	}
}

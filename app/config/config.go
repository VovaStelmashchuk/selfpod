package config

import "os"

var (
	ACastHookToken = os.Getenv("A_CAST_HOOK_TOKEN")

	ACastShowId = os.Getenv("A_CAST_SHOW_ID")

	Port = os.Getenv("APP_PORT")

	GoogleRedirectHost = os.Getenv("GOOGLE_REDIRECT_URL")

	GoogleRedirectPath = "/auth/google/callback"

	YoutubeChannelId = os.Getenv("YOUTUBE_CHANNEL_ID")
)

package config

import "os"

var (
	ACastHookToken     = os.Getenv("A_CAST_HOOK_TOKEN")
	Port               = os.Getenv("APP_PORT")
	GoogleRedirectHost = os.Getenv("GOOGLE_REDIRECT_URL")
	GoogleRedirectPath = "/auth/google/callback"
)

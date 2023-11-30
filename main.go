package main

import (
	"errors"
	"fmt"
	"main/app"
	"main/app/config"
	"main/app/youtube_uploader"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/acast", app.AcastWebHook)
	http.HandleFunc("/login", youtube_uploader.CreateClientEndPoint)
	http.HandleFunc(config.GoogleRedirectPath, youtube_uploader.Oauth2Callback)

	err := http.ListenAndServe(":"+config.Port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

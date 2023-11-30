package main

import (
	"errors"
	"fmt"
	"main/app"
	"main/app/youtube_uploader"
	"net/http"
	"os"
)

const serverPort = "5000"

func main() {
	/*youtube_uploader.UploadToYoutube(
		youtube_uploader.YoutubeUploadRequset{
			Filename:    "tmp_files/output.mov",
			Title:       "test_new_1",
			Description: "test_description_new_1",
		},
	)*/

	http.HandleFunc("/acast", app.AcastWebHook)
	http.HandleFunc("/login", youtube_uploader.CreateClientEndPoint)
	http.HandleFunc("/oauth2callback", youtube_uploader.Oauth2Callback)

	err := http.ListenAndServe(":"+serverPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"errors"
	"fmt"
	"main/app/acast"
	"main/app/config"
	"main/app/google"
	"net/http"
	"os"
)

func main() {
	startServer()
}

func startServer() {
	http.HandleFunc("/acast", acast.WebHook)
	http.HandleFunc("/login", google.LoginToGoogle)
	http.HandleFunc(config.GoogleRedirectPath, google.Oauth2Callback)

	err := http.ListenAndServe(":"+config.Port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

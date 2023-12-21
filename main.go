package main

import (
	"log"
	"main/app/acast"
	"main/app/config"
	"main/app/google"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.OpenFile("files/selfpod-logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	startServer()
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/acast", acast.WebHook)
	mux.HandleFunc("/login", google.LoginToGoogle)
	mux.HandleFunc(config.GoogleRedirectPath, google.Oauth2Callback)

	log.Println("Server starting...")
	err := http.ListenAndServe(":"+config.Port, mux)

	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %s %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

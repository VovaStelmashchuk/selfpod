package main

import (
	"bytes"
	"io"
	"log"
	"main/app/acast"
	"main/app/config"
	"main/app/google"
	"net/http"
	"os"
	"strings"
)

func main() {
	logFile, err := os.OpenFile("selfpod-logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	loggedMux := requestLoggingMiddleware(mux)

	log.Println("Server starting...")
	err := http.ListenAndServe(":"+config.Port, loggedMux)

	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
}

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
			}
			defer r.Body.Close()

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			headers := make([]string, 0, len(r.Header))
			for k, v := range r.Header {
				headers = append(headers, k+": "+strings.Join(v, ", "))
			}

			log.Printf(
				"Received request: %s %s from %s\nHeaders: %s\nBody: %s",
				r.Method, r.URL.Path, r.RemoteAddr, strings.Join(headers, "; "), string(bodyBytes),
			)

			next.ServeHTTP(w, r)
		},
	)
}

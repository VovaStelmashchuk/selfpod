package main

import (
	"errors"
	"fmt"
	"main/app"
	"net/http"
	"os"
)

const serverPort = "5000"

func main() {
	http.HandleFunc("/acast", app.AcastWebHook)
	err := http.ListenAndServe(":"+serverPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

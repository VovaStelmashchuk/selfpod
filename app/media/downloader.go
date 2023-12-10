package media

import (
	"io"
	"log"
	"net/http"
	"os"
)

func downloadAndSaveFile(url string, fileName string) {
	log.Printf("Downloading file %s", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("File %s downloaded successfully into %s", url, fileName)
}

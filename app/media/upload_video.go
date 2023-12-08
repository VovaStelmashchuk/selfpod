package media

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"main/app/google"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"
)

var (
	category = flag.String("category", "22", "Video category")
	keywords = flag.String("keywords", "", "Comma separated list of video keywords")
	privacy  = flag.String("privacy", "unlisted", "Video privacy status")
)

type YoutubeUploadRequset struct {
	Filename    string
	Title       string
	Description string
}

func UploadToYoutube(uploadRequest YoutubeUploadRequset) {
	flag.Parse()

	client, err := google.GetClient(youtube.YoutubeUploadScope)

	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       uploadRequest.Title,
			Description: uploadRequest.Description,
			CategoryId:  *category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: *privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(*keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(*keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(uploadRequest.Filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", uploadRequest.Filename, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		fmt.Printf("Error making YouTube API call: %v", err)
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}

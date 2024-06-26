package media

import (
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"main/app/config"
	"main/app/google"
	"main/app/notifications"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"
)

var (
	category = flag.String("category", "22", "Video category")
	keywords = flag.String("keywords", "", "Comma separated list of video keywords")
	privacy  = flag.String("privacy", "public", "Video privacy status")
)

type YoutubeUploadRequest struct {
	Filename    string
	Title       string
	Description string
}

func UploadToYoutube(uploadRequest YoutubeUploadRequest) error {
	flag.Parse()

	client, getClientError := google.GetClient(youtube.YoutubeUploadScope)

	if getClientError != nil {
		notifications.SendDiscordNotification("Error creating YouTube client, check logs" + getClientError.Error())
		log.Printf("Error creating YouTube client: %v", getClientError)
		return getClientError
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Printf("Error creating YouTube new service: %v", err)
		return err
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       uploadRequest.Title,
			Description: uploadRequest.Description,
			CategoryId:  *category,
			ChannelId:   config.YoutubeChannelId,
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
		notifications.SendDiscordNotification("Error making YouTube API call, check logs")
		log.Printf("Error making YouTube API call: %v", err)
	}
	log.Printf("Upload video result: %v", response)
	return err
}

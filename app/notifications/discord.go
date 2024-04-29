package notifications

import (
	"github.com/gtuk/discordwebhook"
	"log"
	"main/app/config"
)

func SendDiscordNotification(content string) {
	log.Printf("Sending Discord notification: %s", content)

	var username = config.DiscordBotName
	var url = config.DiscordWebhookUrl

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(url, message)
	if err != nil {
		log.Fatal(err)
	}
}

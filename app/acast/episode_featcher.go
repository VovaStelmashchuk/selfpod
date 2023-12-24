package acast

import (
	"bytes"
	"errors"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"log"
	"main/app/config"
	"strings"
	"time"
)

type EpisodeMetaInfo struct {
	Title       string
	Description string
}

func GetEpisodeMetaInfo(episodeId string) (EpisodeMetaInfo, error) {
	log.Printf("Before fetching feed\n")
	time.Sleep(10 * time.Second)
	log.Printf("After fetching feed\n")

	fp := gofeed.NewParser()
	feedUrl := "https://feeds.acast.com/public/shows/" + config.ACastShowId
	log.Printf("Fetching feed from %s\n", feedUrl)
	feed, _ := fp.ParseURL(feedUrl)

	log.Printf("Feed : %s\n", feed)

	found := false
	var targetEpisode *gofeed.Item

	for _, item := range feed.Items {
		if acastEpisodeID := item.Extensions["acast"]["episodeId"][0].Value; acastEpisodeID != "" {
			if acastEpisodeID == episodeId {
				targetEpisode = item
				found = true
				break
			}
		}
	}

	if found {
		return buildEpisodeMetaInfo(targetEpisode)
	} else {
		return EpisodeMetaInfo{}, errors.New("episode not found")
	}
}

func buildEpisodeMetaInfo(targetEpisode *gofeed.Item) (EpisodeMetaInfo, error) {
	lines, err := extractPTagsBeforeBR(targetEpisode.Description)
	return EpisodeMetaInfo{
		Title:       targetEpisode.Title,
		Description: strings.Join(lines, "\n"),
	}, err
}
func extractPTagsBeforeBR(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var pContents []string
	var stop bool

	var f func(*html.Node)
	f = func(n *html.Node) {
		if stop {
			return
		}

		if n.Type == html.ElementNode {
			if n.Data == "br" {
				stop = true
				return
			}

			if n.Data == "p" {
				var buf bytes.Buffer
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						buf.WriteString(c.Data)
					}
				}
				pContents = append(pContents, buf.String())
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for i := range pContents {
		pContents[i] = strings.ReplaceAll(pContents[i], " ", " ")
	}

	return pContents, nil
}

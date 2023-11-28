package app

import (
	"bytes"
	"fmt"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"strings"
)

type YoutubeEpisode struct {
	Title       string
	Description string
}

func GetEpisodeMetaInfo(episodeId string) YoutubeEpisode {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://feeds.acast.com/public/shows/62efce09bcb3d10013e2cc9b")

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
		return buildYoutubeEpisode(targetEpisode)
	} else {
		panic(fmt.Sprintf("Episode with ID %s not found", episodeId))
	}
}

func buildYoutubeEpisode(targetEpisode *gofeed.Item) YoutubeEpisode {
	lines, err := extractPTagsBeforeBR(targetEpisode.Description)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing episode description: %s", err))
	}
	return YoutubeEpisode{
		Title:       targetEpisode.Title,
		Description: strings.Join(lines, "\n"),
	}
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

	return pContents, nil
}

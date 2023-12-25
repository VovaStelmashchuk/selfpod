package processing

import (
	"bytes"
	"errors"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"log"
	"main/app/config"
	"net/http"
	"strings"
	"time"
)

func GetEpisodeDescription(episodeId string) (string, error) {
	feedUrl := "https://feeds.acast.com/public/shows/" + config.ACastShowId
	log.Printf("Fetching feed from %s\n", feedUrl)
	feed, _ := fetchFeed(feedUrl)

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
		return buildEDescription(targetEpisode)
	} else {
		return "", errors.New("episode not found")
	}
}

func fetchFeed(feedUrl string) (*gofeed.Feed, error) {
	// Create a client with some timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Create a request and modify its headers
	req, err := http.NewRequest("GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cache-Control", "no-cache")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the feed
	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func buildEDescription(targetEpisode *gofeed.Item) (string, error) {
	lines, err := extractPTagsBeforeBR(targetEpisode.Description)
	return strings.Join(lines, "\n"), err
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

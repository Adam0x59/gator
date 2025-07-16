package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	feedData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feedStruct RSSFeed
	err = xml.Unmarshal(feedData, &feedStruct)
	if err != nil {
		return nil, err
	}

	feedStruct.Channel.Title = html.UnescapeString(feedStruct.Channel.Title)
	feedStruct.Channel.Description = html.UnescapeString(feedStruct.Channel.Description)

	for i := range feedStruct.Channel.Item {
		feedStruct.Channel.Item[i].Title = html.UnescapeString(feedStruct.Channel.Item[i].Title)
		feedStruct.Channel.Item[i].Description = html.UnescapeString(feedStruct.Channel.Item[i].Description)
	}

	return &feedStruct, nil
}

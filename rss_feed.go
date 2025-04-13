package main

import (
	"context"
	"encoding/xml"
	"errors"
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

func CleanString(feedString string) string {
	newString := html.UnescapeString(feedString)
	return newString
}

func (item *RSSItem) CleanItems() {
	item.Title = CleanString(item.Title)
	item.Link = CleanString(item.Link)
	item.Description = CleanString(item.Description)
	item.PubDate = CleanString(item.PubDate)
}

func (feed *RSSFeed) CleanItems() {
	feed.Channel.Title = CleanString(feed.Channel.Title)
	feed.Channel.Link = CleanString(feed.Channel.Link)
	feed.Channel.Description = CleanString(feed.Channel.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return &RSSFeed{}, errors.New("unable to create request")
	}
	httpClient := http.Client{}
	results, reqError := httpClient.Do(req)

	if reqError != nil {
		return &RSSFeed{}, errors.New("unable to make http request")
	}

	resData, ioError := io.ReadAll(results.Body)

	if ioError != nil {
		return &RSSFeed{}, errors.New("unable to read body from http response")
	}

	xmlData := RSSFeed{}
	umError := xml.Unmarshal(resData, &xmlData)

	if umError != nil {
		return &RSSFeed{}, errors.New("unable to unmarshal data")
	}

	return &xmlData, nil
}
func (feed *RSSFeed) CleanFeed() {
	feed.CleanItems()

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].CleanItems()
	}
}

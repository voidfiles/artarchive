package main

import (
	"log"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestFindImageUrls(t *testing.T) {
	var testTable = []struct {
		content string
		output  []string
	}{
		{"<html><body><div><a></div></html>", []string{}},
		{"<img src='http://example.org'>", []string{"http://example.org"}},
	}

	for _, tt := range testTable {
		t.Run(tt.content, func(t *testing.T) {
			s := FindImageUrls(tt.content)
			assert.Equal(t, s, tt.output)
		})
	}
}

func TestBestContent(t *testing.T) {
	var testTable = []struct {
		inA string
		inB string
		out string
	}{
		{"a", "b", "a"},
		{"aa", "b", "aa"},
		{"aa", "bbb", "bbb"},
	}

	for _, tt := range testTable {
		t.Run(tt.inA, func(t *testing.T) {
			s := BestContent(tt.inA, tt.inB)
			assert.Equal(t, s, tt.out)
		})
	}
}

func TestResolveImageURLs(t *testing.T) {
	var testTable = []struct {
		inBaseLink  string
		inBaseURLs  []string
		outBaseURLs []string
	}{
		{"http://google.com", []string{"b"}, []string{"http://google.com/b"}},
		{"c", []string{"b"}, []string{"/b"}},
		{"http://google.com", []string{"http://googled.com"}, []string{"http://googled.com"}},
	}

	for _, tt := range testTable {
		log.Print(tt)
		s := ResolveImageURLs(tt.inBaseLink, tt.inBaseURLs)
		log.Print(s)
		assert.Equal(t, s[0], tt.outBaseURLs[0])
	}
}

func TestSlidesFromFeeditem(t *testing.T) {
	var testTable = []struct {
		content    string
		slideCount int
	}{
		{"<html><body><div><a></div></html>", 0},
		{"<img src='http://example.org'>", 1},
	}

	for _, tt := range testTable {
		unparsedTime := "2018-10-11T12:00.11Z"
		parsedTime, _ := time.Parse(time.RFC3339Nano, unparsedTime)
		item := gofeed.Item{
			Title:           "A Title",
			Description:     "A description",
			Content:         tt.content,
			Link:            "http://example.org/1",
			Updated:         unparsedTime,
			UpdatedParsed:   &parsedTime,
			Published:       unparsedTime,
			PublishedParsed: &parsedTime,
			GUID:            "http://example.org/1",
		}
		items := make([]*gofeed.Item, 1)
		items[0] = &item
		feed := gofeed.Feed{
			Title:           "A Feed",
			Description:     "A feed description",
			Link:            "http://example.org",
			FeedLink:        "http://example.org/rss",
			Updated:         unparsedTime,
			UpdatedParsed:   &parsedTime,
			Published:       unparsedTime,
			PublishedParsed: &parsedTime,
			Items:           items,
		}
		log.Print(tt)
		slides := SlidesFromFeeditem(&item, feed)
		log.Print(slides)
		assert.Len(t, slides, tt.slideCount)

	}
}

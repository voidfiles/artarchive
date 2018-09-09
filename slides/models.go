package slides

import (
	"time"
)

type Binding struct {
	In  chan Slide
	Out chan Slide
}

// Site contains identifying info about what parent site an image was found
type Site struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Page contains identifying info about where an image was found
type Page struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Published *time.Time `json:"published"`
	GUIDHash  string
}

// ImageInfo is important information about a single image
type ImageInfo struct {
	URL         string `json:"url,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

// ArtistInfo is important information about the artist who made the image
type ArtistInfo struct {
	Name         string   `json:"name,omitempty"`
	ArtsyURL     string   `json:"artsy_url,omitempty"`
	WikipediaURL string   `json:"wikipedia_url,omitempty"`
	WebsiteURL   string   `json:"website_url,omitempty"`
	FeedURL      string   `json:"feed_url,omitempty"`
	InstagramURL string   `json:"instagram_url,omitempty"`
	TwitterURL   string   `json:"twitter_url,omitempty"`
	Description  string   `json:"description,omitempty"`
	Feeds        []string `json:"feeds,omitempty"`
	Sites        []string `json:"sites,omitempty"`
}

// WorkInfo is important information about the work
type WorkInfo struct {
	Name string `json:"name,omitempty"`
}

// Slide bundles together important information about a find
type Slide struct {
	Site           Site          `json:"site,omitempty"`
	Page           Page          `json:"page,omitempty"`
	Content        string        `json:"content,omitempty"`
	GUIDHash       string        `json:"guid_hash,omitempty"`
	SourceImageURL string        `json:"source_image_url,omitempty"`
	ArchivedImage  *ImageInfo    `json:"archived_image,omitempty"`
	ArtistInfo     *ArtistInfo   `json:"artist_info,omitempty"`
	ArtistsInfo    []*ArtistInfo `json:"artists,omitempty"`
	WorkInfo       *WorkInfo     `json:"work_info,omitempty"`
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImageURL(t *testing.T) {
	testTable := []struct {
		inURL  string
		outURL string
	}{
		{
			"https://78.media.tumblr.com/4d27b79856711ef14c949c1bc89e5457/tumblr_p5jq9r2xhj1qz9v0to2_500.jpg",
			"https://78.media.tumblr.com/4d27b79856711ef14c949c1bc89e5457/tumblr_p5jq9r2xhj1qz9v0to2_1280.jpg",
		}, {
			"http://78.media.tumblr.com/efaf405a1476afdd03ac26c24717d228/tumblr_inline_p5hd3jkleZ1se9q8s_500.jpg",
			"https://78.media.tumblr.com/efaf405a1476afdd03ac26c24717d228/tumblr_inline_p5hd3jkleZ1se9q8s_1280.jpg",
		}, {
			"https://78.media.tumblr.com/2af2523f8d048c568c01b7f616a0b9eb/tumblr_inline_p56v0jZq1t1qz9p2o_540.jpg",
			"https://78.media.tumblr.com/2af2523f8d048c568c01b7f616a0b9eb/tumblr_inline_p56v0jZq1t1qz9p2o_1280.jpg",
		},
	}
	for _, test := range testTable {
		assert.Equal(t, test.outURL, ParseImageURL(test.inURL))
	}

}

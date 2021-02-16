package types

import "time"

type Link struct {
	Subreddit   string    `json:"subreddit"`
	Content     string    `json:"selftext"`
	Title       string    `json:"title"`
	UpvoteRatio float64   `json:"upvote_ratio"`
	Upvotes     int       `json:"ups"`
	Created     time.Time `json:"created"`
	CreatedUtc  time.Time `json:"created_utc"`
	HtmlContent string    `json:"selftext_html"`
	NSFW        bool      `json:"over_18"`
	Spoiler     bool      `json:"spoiler"`
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Permalink   string    `json:"permalink"`
}

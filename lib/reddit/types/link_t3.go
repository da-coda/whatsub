package types

type Link struct {
	Kind string
	Data struct {
		Modhash  string
		Dist     int
		Children []Post
	}
}

type Post struct {
	Kind string
	Data struct {
		Subreddit   string     `json:"subreddit_name_prefixed"`
		Content     string     `json:"selftext"`
		Title       string     `json:"title"`
		UpvoteRatio float64    `json:"upvote_ratio"`
		Upvotes     int        `json:"ups"`
		Created     RedditTime `json:"created"`
		CreatedUtc  RedditTime `json:"created_utc"`
		HtmlContent string     `json:"selftext_html"`
		NSFW        bool       `json:"over_18"`
		Spoiler     bool       `json:"spoiler"`
		ID          string     `json:"id"`
		Author      string     `json:"author"`
		Permalink   string     `json:"permalink"`
		Image       string     `json:"url"`
	}
}

func (l Link) GetContent() []Post {
	return l.Data.Children
}

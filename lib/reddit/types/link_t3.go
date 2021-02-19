package types

type Link struct {
	Kind string
	Data struct {
		Modhash  string
		Dist     int
		Children []Post
	}
}

type PostType string

const (
	TextPost  PostType = "Text"
	ImagePost PostType = "Image"
	LinkPost  PostType = "Link"
	VideoPost PostType = "Video"
)

type Post struct {
	Kind string
	Data struct {
		Subreddit       string     `json:"subreddit_name_prefixed"`
		Content         string     `json:"selftext"`
		Title           string     `json:"title"`
		UpvoteRatio     float64    `json:"upvote_ratio"`
		Upvotes         int        `json:"ups"`
		Created         RedditTime `json:"created"`
		CreatedUtc      RedditTime `json:"created_utc"`
		HtmlContent     string     `json:"selftext_html"`
		NSFW            bool       `json:"over_18"`
		Spoiler         bool       `json:"spoiler"`
		ID              string     `json:"id"`
		Author          string     `json:"author"`
		Permalink       string     `json:"permalink"`
		Url             string     `json:"url"`
		IsSelf          bool       `json:"is_self"`
		IsVideo         bool       `json:"is_video"`
		PostHint        string     `json:"post_hint"`
		CrossPostParent string     `json:"crosspost_parent,omitempty"`
	}
}

func (l Link) GetContent() []Post {
	return l.Data.Children
}

func (p Post) GetType() PostType {
	if p.Data.IsSelf {
		return TextPost
	}
	if p.Data.IsVideo {
		return VideoPost
	}
	if p.Data.PostHint == "image" {
		return ImagePost
	}
	if p.Data.Url != "" {
		return LinkPost
	}
	return ""
}

func (p Post) IsCrosspost() bool {
	return p.Data.CrossPostParent != ""
}

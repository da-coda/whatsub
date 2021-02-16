package redditHelper

import (
	"github.com/da-coda/whatsub/lib/reddit"
	"github.com/da-coda/whatsub/lib/reddit/types"
)

func GetTopPostsForSubreddits(subreddits []string, amountPerSubreddit int) ([]types.Link, error) {
	var allPosts []types.Link
	for _, sub := range subreddits {
		posts, err := reddit.SubredditTopPosts(sub, amountPerSubreddit, 0)
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, posts)
	}
	return allPosts, nil
}

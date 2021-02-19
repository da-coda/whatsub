package reddit

import (
	"encoding/json"
	"fmt"
	"github.com/da-coda/whatsub/lib/reddit/types"
	"github.com/pkg/errors"
	"strconv"
)

func SubredditTopPosts(subreddit string, limit int, count int) (types.Link, error) {
	uri := fmt.Sprintf("%s/%s/top.json", baseUrl, subreddit)
	response, err := client.R().SetQueryParams(map[string]string{
		"limit": strconv.Itoa(limit),
		"count": strconv.Itoa(count),
	}).Get(uri)
	if err != nil {
		return types.Link{}, errors.Wrapf(err, "Unable to get top posts for subreddit %s", subreddit)
	}
	var posts types.Link
	err = json.Unmarshal(response.Body(), &posts)
	if err != nil {
		return types.Link{}, errors.Wrapf(err, "Unable to marshal posts for subreddit %s", subreddit)
	}
	return posts, nil
}

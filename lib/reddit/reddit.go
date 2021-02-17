package reddit

import (
	_ "encoding/json"
	"gopkg.in/resty.v1"
)

const baseUrl = "https://reddit.com"

var client *resty.Client

func init() {
	client = resty.New()
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))
}

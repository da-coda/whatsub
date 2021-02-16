package main

import "gopkg.in/resty.v1"

const baseUrl = "https://reddit.com/"

var client *resty.Client

func init() {
	client = resty.New()
}

package messages

import "github.com/da-coda/whatsub/lib/reddit/types"

const RoundType Type = "round"

type Round struct {
	Envelop
	Payload struct {
		Number int
		From   int
		Post   struct {
			Title   string
			Content string
			Type    types.PostType
			Url     string
		}
		Subreddits []string
	}
}

func NewRoundMessage() Round {
	return Round{
		Envelop: Envelop{Type: RoundType},
		Payload: struct {
			Number int
			From   int
			Post   struct {
				Title   string
				Content string
				Type    types.PostType
				Url     string
			}
			Subreddits []string
		}{},
	}
}

package types

import (
	"strconv"
	"strings"
	"time"
)

type RedditTime struct {
	Time time.Time
}

func (rt *RedditTime) UnmarshalJSON(data []byte) error {
	timestamp := string(data)
	truncatedData := strings.Split(timestamp, ".")[0]
	q, err := strconv.ParseInt(truncatedData, 10, 64)
	if err != nil {
		return err
	}
	rt.Time = time.Unix(q, 0)
	return nil
}

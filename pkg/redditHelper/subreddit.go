package redditHelper

import (
	"context"
	"github.com/da-coda/whatsub/pkg/database"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	TopSubreddits = `SELECT * FROM subreddits ORDER BY subscriber_count DESC LIMIT $1`
)

type Subreddit struct {
	Name     string    `db:"sub_name"`
	Path     string    `db:"sub_full_name"`
	SubCount int       `db:"subscriber_count"`
	Created  time.Time `db:"created_date"`
	Updated  time.Time `db:"last_updated"`
}

type Subreddits []Subreddit

func GetTopSubreddits(limit int) (Subreddits, error) {
	var subreddits Subreddits
	db := database.GetConn()
	ctx := context.Background()
	timeout, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	txx, err := db.BeginTxx(timeout, nil)
	if err != nil {
		cancelFunc()
		return Subreddits{}, errors.Wrap(err, "Unable to start transaction")
	}
	defer func() {
		err := txx.Commit()
		cancelFunc()
		if err != nil {
			logrus.WithError(err).Error("Unable to close transaction")
		}
	}()
	rows, err := txx.Queryx(TopSubreddits, limit)
	if err != nil {
		return Subreddits{}, errors.Wrap(err, "Unable to query db")
	}
	for rows.Next() {
		var sub Subreddit
		err = rows.StructScan(&sub)
		if err != nil {
			return Subreddits{}, errors.Wrap(err, "Unable to scan row")
		}
		subreddits = append(subreddits, sub)
	}
	return subreddits, nil
}

func (subs Subreddits) GetSubPaths() []string {
	var paths []string
	for _, sub := range subs {
		paths = append(paths, sub.Path)
	}
	return paths
}

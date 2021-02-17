package worker

import (
	"fmt"
	"github.com/da-coda/whatsub/lib/reddit/types"
	"github.com/da-coda/whatsub/pkg/redditHelper"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const InactiveNotStartedTimeout = 60 * time.Minute

type Worker struct {
	WorkerId         uuid.UUID
	Connections      map[string]*websocket.Conn
	connectionLookup map[*websocket.Conn]string
	Started          bool
	Posts            []types.Post
	Subreddits       []string
	Host             string
	Rounds           int
	Active           bool
	Created          time.Time
}

type Round struct {
	Type      string
	Number    int
	From      int
	PostTitle string
	PostText  string
	Image     string
}

func New() *Worker {
	w := &Worker{WorkerId: uuid.New(), Rounds: 10, Created: time.Now(), Started: false, Active: false}
	w.Connections = make(map[string]*websocket.Conn)
	w.connectionLookup = make(map[*websocket.Conn]string)
	return w
}

func (worker *Worker) AddPlayer(conn *websocket.Conn) bool {
	if worker.Started {
		return false
	}
	_, message, err := conn.ReadMessage()
	fmt.Print(string(message))
	if worker.Host == "" {
		worker.Host = string(message)
	}
	worker.Connections[string(message)] = conn
	worker.connectionLookup[conn] = string(message)
	logrus.WithField("Worker", worker.WorkerId).Info("New player joined")
	err = conn.WriteJSON(`{"message": "welcome ` + string(message) + `"}`)
	if err != nil {
		logrus.WithError(err).WithField("Worker", worker.WorkerId).Error("Unable to greet new player")
		return false
	}

	return true
}

func (worker *Worker) RunGame() {
	worker.Started = true
	worker.Active = true
	worker.preparePosts()
	for i := 0; i < worker.Rounds; i++ {
		worker.runRound(i)
	}
	worker.Active = false
}

func (worker *Worker) preparePosts() {
	subreddits := redditHelper.GetTopSubreddits()
	links, err := redditHelper.GetTopPostsForSubreddits(subreddits, 5)
	if err != nil {
		logrus.WithError(err).Error("Unable to prepare posts")
	}
	var posts []types.Post
	for _, link := range links {
		for _, linkPost := range link.GetContent() {
			posts = append(posts, linkPost)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })
	worker.Posts = posts
	worker.Subreddits = subreddits
}

func (worker *Worker) runRound(round int) {
	var wg sync.WaitGroup
	post := worker.Posts[round]
	postText := post.Data.HtmlContent
	if postText == "" {
		postText = post.Data.Content
	}
	roundPayload := Round{
		Type:      "round",
		Number:    round,
		From:      worker.Rounds,
		PostTitle: post.Data.Title,
		PostText:  postText,
		Image:     post.Data.Image,
	}
	logrus.Debug(post.Data.Subreddit)
	for _, client := range worker.Connections {
		err := client.WriteJSON(roundPayload)
		if err != nil {
			logrus.WithError(err).Error("Unable to ping player")
		}
		wg.Add(1)
		go handleClientAnswer(client, post.Data.Subreddit, &wg)
	}
	wg.Wait()
}

// StillNeeded checks for different conditions to decide if this worker is still needed
func (worker *Worker) StillNeeded() bool {
	// All rounds are played, game is done
	if !worker.Active && worker.Started {
		logrus.Debug("Game is done")
		return false
	}

	// Worker got created but game didn't start within the duration InactiveNotStartedTimeout
	if worker.Created.Add(InactiveNotStartedTimeout).Before(time.Now()) && !worker.Started {
		logrus.Debug("Game never started")
		return false
	}

	return true
}

func handleClientAnswer(client *websocket.Conn, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, answer, err := client.ReadMessage()
	if err != nil {
		logrus.WithError(err).Error("Unable to get answer from client")
		return
	}
	if strings.Compare(string(answer), correctAnswer) != 0 {
		err := client.WriteMessage(websocket.TextMessage, []byte("false"))
		if err != nil {
			logrus.WithError(err).Error("Unable to notify client about answer correctness")
		}
		return
	}
	err = client.WriteMessage(websocket.TextMessage, []byte("true"))
	if err != nil {
		logrus.WithError(err).Error("Unable to notify client about answer correctness")
	}
}

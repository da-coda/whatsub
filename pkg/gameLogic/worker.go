package gameLogic

import (
	"github.com/da-coda/whatsub/lib/reddit/types"
	"github.com/da-coda/whatsub/messages"
	"github.com/da-coda/whatsub/pkg/redditHelper"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const InactiveNotStartedTimeout = 60 * time.Minute

type Worker struct {
	WorkerId    uuid.UUID
	Connections map[string]*Client
	Started     bool
	Posts       []types.Post
	Subreddits  []string
	Host        string
	Rounds      int
	Active      bool
	Created     time.Time
	Register    chan *Client
	Unregister  chan *Client
}

func NewWorker() *Worker {
	w := &Worker{
		WorkerId:    uuid.New(),
		Rounds:      10,
		Created:     time.Now(),
		Started:     false,
		Active:      false,
		Register:    make(chan *Client),
		Connections: make(map[string]*Client),
	}
	return w
}

func (worker *Worker) OpenLobby() {
	logrus.WithField("Worker", worker.WorkerId).Debug("Open Lobby")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for !worker.Started {
		select {
		case gameClient := <-worker.Register:
			logrus.WithField("Worker", worker.WorkerId).WithField("Player", gameClient.Name).Debug("Player joined")
			worker.Connections[gameClient.Name] = gameClient
		case <-ticker.C:
			continue
		}
	}
	logrus.WithField("Worker", worker.WorkerId).Debug("Closing Lobby")
}

func (worker *Worker) RunGame() {
	worker.Started = true
	worker.Active = true
	logrus.WithField("Worker", worker.WorkerId).Debug("Starting Game")
	worker.preparePosts()
	for i := 0; i < worker.Rounds; i++ {
		worker.runRound(i)
	}
	for _, conn := range worker.Connections {
		_ = conn.WriteJSON(map[string]string{"Score": strconv.Itoa(conn.Score)})
	}
	worker.Active = false
}

func (worker *Worker) preparePosts() {
	logrus.WithField("Worker", worker.WorkerId).Debug("Preparing Posts")
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
	logrus.WithField("Worker", worker.WorkerId).WithField("Round", round).Debug("Starting Round")
	var wg sync.WaitGroup
	post := worker.Posts[round]
	postText := post.Data.HtmlContent
	if postText == "" {
		postText = post.Data.Content
	}
	roundMessage := messages.NewRoundMessage()
	roundMessage.Payload.Number = round
	roundMessage.Payload.From = worker.Rounds
	roundMessage.Payload.Post.Title = post.Data.Title
	roundMessage.Payload.Post.Content = postText
	roundMessage.Payload.Post.Type = post.GetType()
	roundMessage.Payload.Post.Url = post.Data.Url

	for _, playerClient := range worker.Connections {
		err := playerClient.WriteJSON(roundMessage)
		if err != nil {
			logrus.WithError(err).Error("Unable to ping player")
		}
		wg.Add(1)
		go worker.handleClientAnswer(playerClient, post.Data.Subreddit, &wg)
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

func (worker *Worker) handleClientAnswer(playerClient *Client, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, answer, err := playerClient.ReadMessage()
	if err != nil {
		logrus.WithError(err).Error("Unable to get answer from client")
		return
	}
	if strings.Compare(string(answer), correctAnswer) != 0 {
		err := playerClient.WriteMessage(websocket.TextMessage, []byte("false"))
		if err != nil {
			logrus.WithError(err).Error("Unable to notify client about answer correctness")
		}
		return
	}
	playerClient.Score++
	err = playerClient.WriteMessage(websocket.TextMessage, []byte("true"))
	if err != nil {
		logrus.WithError(err).Error("Unable to notify client about answer correctness")
	}
}

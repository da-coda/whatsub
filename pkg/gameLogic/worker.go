package gameLogic

import (
	"encoding/json"
	"github.com/da-coda/whatsub/lib/reddit/types"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/da-coda/whatsub/pkg/redditHelper"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const InactiveNotStartedTimeout = 60 * time.Minute

type State int

const (
	Created State = iota
	Open
	Started
	Done
)

type Worker struct {
	Id          uuid.UUID
	Connections map[string]*Client
	Posts       []types.Post
	Subreddits  []string
	Host        string
	Rounds      int
	Created     time.Time
	Register    chan *Client
	Unregister  chan *Client
	State       State
	Incoming    chan *Client
}

func NewWorker() *Worker {
	w := &Worker{
		Id:          uuid.New(),
		Rounds:      10,
		Created:     time.Now(),
		State:       Created,
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Incoming:    make(chan *Client),
		Connections: make(map[string]*Client),
	}
	return w
}

func (worker *Worker) OpenLobby() {
	logrus.WithField("Worker", worker.Id).Debug("Open Lobby")
	worker.State = Open
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for worker.State == Open {
		select {
		case gameClient := <-worker.Register:
			logrus.WithField("Worker", worker.Id).WithField("Player", gameClient.Name).Debug("Player joined")
			worker.Connections[gameClient.Name] = gameClient
		case <-ticker.C:
			continue
		}
	}
	logrus.WithField("Worker", worker.Id).Debug("Closing Lobby")
}

func (worker *Worker) DisconnectHandler() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for worker.State >= Open {
		select {
		case gameClient := <-worker.Unregister:
			logrus.WithField("Worker", worker.Id).WithField("Player", gameClient.Name).Debug("Player disconnected")
			delete(worker.Connections, gameClient.Name)
		case <-ticker.C:
			continue
		}
	}
}

func (worker *Worker) RunGame() {
	worker.State = Started
	logrus.WithField("Worker", worker.Id).Debug("Starting Game")
	worker.preparePosts()
	for i := 0; i < worker.Rounds; i++ {
		worker.runRound(i)
	}
	for _, conn := range worker.Connections {
		msg := messages.NewScoreMessage()
		msg.Payload.Score = conn.Score
		msgJson, err := json.Marshal(msg)
		if err != nil {
			logrus.WithError(err).
				WithField("Worker", worker.Id).
				WithField("Client", conn.Name).
				Error("Unable to marshal score message to json")
			continue
		}
		conn.Send <- msgJson
	}
	worker.State = Done
}

func (worker *Worker) preparePosts() {
	logrus.WithField("Worker", worker.Id).Debug("Preparing Posts")
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
	logrus.WithField("Worker", worker.Id).WithField("Round", round).Info("Starting Round")
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
	roundJson, err := json.Marshal(roundMessage)
	if err != nil {
		logrus.WithError(err).
			WithField("Worker", worker.Id).
			Error("Unable to marshal round message to json")
		return
	}
	for _, playerClient := range worker.Connections {
		playerClient.Send <- roundJson
	}
	for i := 0; i < len(worker.Connections); i++ {
		clientAnswered := <-worker.Incoming
		logrus.WithField("Worker", worker.Id).WithField("Client", clientAnswered.Name).Info("Client answered")
		wg.Add(1)
		go worker.handleClientAnswer(clientAnswered, post.Data.Subreddit, &wg)
	}

	wg.Wait()
}

// StillNeeded checks for different conditions to decide if this worker is still needed
func (worker *Worker) StillNeeded() bool {
	// All rounds are played, game is done
	if worker.State == Done {
		logrus.Debug("Game is done")
		return false
	}

	// Worker got created but game didn't start within the duration InactiveNotStartedTimeout
	if worker.State == Created && worker.Created.Add(InactiveNotStartedTimeout).Before(time.Now()) {
		logrus.Debug("Game never started")
		return false
	}

	return true
}

func (worker *Worker) handleClientAnswer(playerClient *Client, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	answer := <-playerClient.Message
	logrus.WithField("Worker", worker.Id).WithField("Client", playerClient.Name).WithField("Answer", answer).Info("Client answered")
	msg := messages.NewAnswerCorrectnessMessage()
	msg.Payload.Correct = strings.Compare(string(answer), correctAnswer) == 0
	if msg.Payload.Correct {
		playerClient.Score++
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).
			WithField("Worker", worker.Id).
			Error("Unable to marshal answer correctness message to json")
		return
	}
	playerClient.Send <- msgJson
}

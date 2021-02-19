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

const InactiveNotStartedTimeout = 1 * time.Hour
const EmptyLobbyTimeout = 10 * time.Minute
const LobbyNotStartedTimeout = 30 * time.Minute

type State int

const (
	Created State = iota
	Open
	Started
	Done
	Closed
)

//Worker handles a single game, holds all participating clients and needed resources for the game
type Worker struct {
	Id          uuid.UUID
	Clients     map[string]*Client
	Posts       []types.Post
	Subreddits  []string
	Host        string
	Rounds      int
	Created     time.Time
	LobbyOpened time.Time
	Register    chan *Client
	Unregister  chan *Client
	State       State
	Incoming    chan *Client
}

//NewWorker creates a new Worker and setups channels
func NewWorker() *Worker {
	w := &Worker{
		Id:         uuid.New(),
		Rounds:     10,
		Created:    time.Now(),
		State:      Created,
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		Incoming:   make(chan *Client, 256),
		Clients:    make(map[string]*Client),
	}
	return w
}

//Close implements the io.Closer interface and closes all channels and calls Client.Close on all connected clients
func (worker *Worker) Close() error {
	worker.State = Closed
	logrus.WithField("Worker", worker.Id).Debug("Terminating worker because Close() got called")
	for _, client := range worker.Clients {
		worker.Unregister <- client
	}
	close(worker.Register)
	close(worker.Incoming)
	return nil
}

//OpenLobby checks if a new worker registers while the worker State is Open
func (worker *Worker) OpenLobby() {
	worker.LobbyOpened = time.Now()
	logrus.WithField("Worker", worker.Id).Debug("Open Lobby")
	worker.State = Open
	go worker.DisconnectHandler()
	//create a ticker and use it in our loop that checks for new registers so that at least every second the loop condition
	//is checked even without register event so that the lobby can be closed correctly
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for worker.State == Open {
		select {
		case gameClient, open := <-worker.Register:
			if open {
				logrus.WithField("Worker", worker.Id).Debug("Player joined")
				worker.Clients[gameClient.Name] = gameClient
			}
		case <-ticker.C:
			continue
		}
	}
	logrus.WithField("Worker", worker.Id).Debug("Closing Lobby")
}

//DisconnectHandler listens on the Unregister channel and removes clients that unregistered themselves
func (worker *Worker) DisconnectHandler() {
	//use ticker so that every second the loop condition is checked even without unregister event
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for worker.State >= Open {
		select {
		case gameClient, open := <-worker.Unregister:
			if open {
				logrus.WithField("Worker", worker.Id).Debug("Player disconnected")
				if gameClient != nil {
					_ = gameClient.Close()
				}
				delete(worker.Clients, gameClient.Name)
			}
		case <-ticker.C:
			continue
		}
	}
}

//RunGame is the main game loop which prepares the posts and runs each round
func (worker *Worker) RunGame() {
	if worker.State == Started {
		return
	}
	worker.State = Started
	logrus.WithField("Worker", worker.Id).Debug("Starting Game")
	worker.preparePosts()

	// run Worker.Rounds rounds
	for i := 0; i < worker.Rounds; i++ {
		worker.runRound(i)
		msg := messages.NewScoreMessage()
		for name, conn := range worker.Clients {
			msg.Payload.Scores[name] = conn.Score
		}
		msgJson, err := json.Marshal(msg)
		if err != nil {
			logrus.WithError(err).
				WithField("Worker", worker.Id).
				Error("Unable to marshal score message to json")
			continue
		}
		for _, conn := range worker.Clients {
			conn.Send <- msgJson
		}
		time.Sleep(2 * time.Second)
	}
	//set State to Done so that the clean up routine of GameMaster can handle the termination of the worker and clients
	worker.State = Done
}

//preparePosts collects subreddits and posts for those subreddits, shuffles them around and adds them to the worker
func (worker *Worker) preparePosts() {
	logrus.WithField("Worker", worker.Id).Debug("Preparing Posts")
	//collect subreddits and posts
	subreddits := redditHelper.GetTopSubreddits()
	links, err := redditHelper.GetTopPostsForSubreddits(subreddits, 5)
	if err != nil {
		logrus.WithError(err).Error("Unable to prepare posts")
	}
	//merge all posts from all subreddits into a single slice
	var posts []types.Post
	for _, link := range links {
		for _, linkPost := range link.GetContent() {
			if linkPost.GetType() != types.LinkPost && linkPost.GetType() != types.VideoPost {
				posts = append(posts, linkPost)
			}
		}
	}
	//shuffle slice randomly and add them to the worker
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })
	worker.Posts = posts
	worker.Subreddits = subreddits
}

//runRound handles a single round by parsing a post into the RoundMessage struct, sending it to all clients and spawning a handler for incoming answers
func (worker *Worker) runRound(round int) {
	logrus.WithField("Worker", worker.Id).WithField("Round", round).Info("Starting Round")
	var wg sync.WaitGroup

	//get the post for this round
	post := worker.Posts[round]
	//if post is an html post use the HtmlContent, otherwise just the Content
	postText := post.Data.HtmlContent
	if postText == "" {
		postText = post.Data.Content
	}

	//Parse the post into our RoundMessage format and marshal it to json
	roundMessage := messages.NewRoundMessage()
	roundMessage.Payload.Number = round
	roundMessage.Payload.From = worker.Rounds
	roundMessage.Payload.Post.Title = post.Data.Title
	roundMessage.Payload.Post.Content = postText
	roundMessage.Payload.Post.Type = post.GetType()
	roundMessage.Payload.Post.Url = post.Data.Url
	roundMessage.Payload.Subreddits = worker.Subreddits
	roundJson, err := json.Marshal(roundMessage)
	if err != nil {
		logrus.WithError(err).
			WithField("Worker", worker.Id).
			Error("Unable to marshal round message to json")
		return
	}

	//send post as json to all clients
	for _, playerClient := range worker.Clients {
		playerClient.Send <- roundJson
		playerClient.Blocked = false
	}

	//listen for incoming answers and spawn a handler for handling the answer
	for i := 0; i < len(worker.Clients); i++ {
		clientAnswered := <-worker.Incoming
		if clientAnswered.Blocked {
			i--
			continue
		}
		worker.handleClientAnswer(clientAnswered, post.Data.Subreddit, &wg)
	}
}

// StillNeeded checks for different conditions to decide if this worker is still needed
func (worker *Worker) StillNeeded() bool {
	// All rounds are played, game is done
	if worker.State == Done {
		logrus.Debug("Game is done")
		return false
	}

	// Lobby been open for EmptyLobbyTimeout minutes without anyone joining
	if worker.State == Open && len(worker.Clients) == 0 && worker.LobbyOpened.Add(EmptyLobbyTimeout).Before(time.Now()) {
		logrus.Debug("Lobby been empty for too long")
		return false
	}

	// Lobby been open for LobbyNotStartedTimeout without starting the game
	if worker.State == Open && worker.LobbyOpened.Add(LobbyNotStartedTimeout).Before(time.Now()) {
		logrus.Debug("Lobby been open for too long")
		return false
	}

	// Worker got created but game didn't start within the duration InactiveNotStartedTimeout
	if worker.State == Created && worker.Created.Add(InactiveNotStartedTimeout).Before(time.Now()) {
		logrus.Debug("Game never started")
		return false
	}

	return true
}

//handleClientAnswer handles the incoming answer of a single client, updates the score if necessary, notifies the client
func (worker *Worker) handleClientAnswer(playerClient *Client, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	//receive answer from client
	answerJson := <-playerClient.Message
	var answerMessage messages.Answer
	err := answerMessage.Parse(answerJson)
	if err != nil {
		logrus.WithError(err).Error("Unable to parse answer message")
		return
	}
	answer := answerMessage.Payload.Answer
	logrus.WithField("Worker", worker.Id).WithField("Client", playerClient.Name).WithField("Answer", answer).Debug("Client answered")

	msg := messages.NewAnswerCorrectnessMessage()
	msg.Payload.Correct = strings.Compare(string(answer), correctAnswer) == 0
	msg.Payload.CorrectAnswer = correctAnswer
	if msg.Payload.Correct {
		playerClient.Score++
	}

	//notify client if correct or not
	msgJson, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).
			WithField("Worker", worker.Id).
			Error("Unable to marshal answer correctness message to json")
		return
	}
	playerClient.Send <- msgJson
	playerClient.Blocked = true
}

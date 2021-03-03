package gameLogic

import (
	"encoding/json"
	"github.com/da-coda/whatsub/lib/reddit/types"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/da-coda/whatsub/pkg/redditHelper"
	"github.com/da-coda/whatsub/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const InactiveNotStartedTimeout = 1 * time.Hour
const LobbyNotStartedTimeout = 30 * time.Minute
const EmptyLobbyTimeout = 10 * time.Minute

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

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
	Id            uuid.UUID
	ShortId       string
	Clients       sync.Map
	Posts         []types.Post
	Score         sync.Map
	Subreddits    []string
	Host          string
	CreatorIpHash string
	RoundsTotal   int
	Created       time.Time
	LobbyOpened   time.Time
	State         State
	ClientCount   uint64
	log           *logrus.Entry
}

//NewWorker creates a new Worker and setups channels
func NewWorker() *Worker {
	w := &Worker{
		Id:          uuid.New(),
		ShortId:     utils.KeyGenerator(8),
		RoundsTotal: 10,
		Created:     time.Now(),
		State:       Created,
		Clients:     sync.Map{},
		ClientCount: 0,
	}
	w.log = logrus.WithField("Worker", w.Id.String())
	return w
}

//Close implements the io.Closer interface and closes all channels and calls Client.Close on all connected clients
func (worker *Worker) Close() error {
	worker.State = Closed
	worker.log.Debug("Terminating worker because Close() got called")
	worker.Clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		worker.Disconnect(client)
		return true
	})
	return nil
}

//OpenLobby checks if a new worker registers while the worker State is Open
func (worker *Worker) Join(w http.ResponseWriter, r *http.Request) {
	playerName := r.FormValue("name")
	playerUUIDString := r.FormValue("uuid")
	playerUUID, err := uuid.Parse(playerUUIDString)
	if err != nil {
		worker.log.WithError(err).Error("invalid uuid")
		w.WriteHeader(400)
		return
	}
	switch worker.State {
	case Started:
		if _, exists := worker.Score.Load(playerUUIDString); exists {
			client, err := worker.join(w, r, playerName, playerUUID)
			if err != nil {
				return
			}
			client.Send <- []byte("Welcome back")
			break
		}
		w.WriteHeader(400)
		return
	case Open:
		client, err := worker.join(w, r, playerName, playerUUID)
		if err != nil {
			worker.log.WithError(err).Error("Dafuq")
			return
		}
		worker.Score.Store(client.uuid.String(), 0)
	default:
		w.WriteHeader(400)
	}
	worker.sendScoreMessage()
}

func (worker *Worker) Disconnect(gameClient *Client) {
	worker.log.Debug("Player disconnected")
	worker.Clients.Delete(gameClient.uuid)
	_ = gameClient.Close()
	worker.ClientCount--
}

//RunGame is the main game loop which prepares the posts and runs each round
func (worker *Worker) RunGame() {
	if worker.State == Started {
		return
	}
	for len(worker.Posts) == 0 {
		time.Sleep(1 * time.Second)
	}
	worker.State = Started
	worker.log.Debug("Starting Game")

	// run Worker.RoundsTotal rounds
	for i := 0; i < worker.RoundsTotal; i++ {
		worker.runRound(i)
		worker.sendScoreMessage()
		time.Sleep(2 * time.Second)
	}
	//set State to Done so that the clean up routine of GameMaster can handle the termination of the worker and clients
	worker.State = Done
}

// StillNeeded checks for different conditions to decide if this worker is still needed
func (worker *Worker) StillNeeded() bool {
	// All rounds are played, game is done
	if worker.State == Done {
		worker.log.Debug("Game is done")
		return false
	}

	// Lobby been open for EmptyLobbyTimeout minutes without anyone joining
	if worker.State == Open && worker.ClientCount == 0 && worker.LobbyOpened.Add(EmptyLobbyTimeout).Before(time.Now()) {
		worker.log.Debug("Lobby been empty for too long")
		return false
	}

	// Lobby been open for LobbyNotStartedTimeout without starting the game
	if worker.State == Open && worker.LobbyOpened.Add(LobbyNotStartedTimeout).Before(time.Now()) {
		worker.log.Debug("Lobby been open for too long")
		return false
	}

	// Worker got created but game didn't start within the duration InactiveNotStartedTimeout
	if worker.State == Created && worker.Created.Add(InactiveNotStartedTimeout).Before(time.Now()) {
		worker.log.Debug("Game never started")
		return false
	}

	// All clients left during the game
	if worker.State == Started && worker.ClientCount == 0 {
		worker.log.Debug("Game abandoned")
		return false
	}

	return true
}

//handleClientAnswer handles the incoming answer of a single client, updates the score if necessary, notifies the client
func (worker *Worker) handleClientAnswer(playerClient *Client, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	//receive answer from client
	answerJson, ok := <-playerClient.Message
	if !ok {
		return
	}
	playerClient.Blocked = true
	var answerMessage messages.Answer
	err := answerMessage.Parse(answerJson)
	if err != nil {
		worker.log.WithError(err).Error("Unable to parse answer message")
		return
	}
	answer := answerMessage.Payload.Answer
	worker.log.WithField("Client", playerClient.Name).WithField("Answer", answer).Trace("Client answered")

	msg := messages.NewAnswerCorrectnessMessage()
	msg.Payload.Correct = strings.Compare(answer, correctAnswer) == 0
	msg.Payload.CorrectAnswer = correctAnswer
	if msg.Payload.Correct {
		currentScore, _ := worker.Score.Load(playerClient.uuid.String())
		worker.Score.Store(playerClient.uuid.String(), 1+currentScore.(int))
	}

	//notify client if correct or not
	msgJson, err := json.Marshal(msg)
	if err != nil {
		worker.log.WithError(err).
			Error("Unable to marshal answer correctness message to json")
		return
	}
	if !playerClient.Terminated {
		playerClient.Send <- msgJson
	}
}

//preparePosts collects subreddits and posts for those subreddits, shuffles them around and adds them to the worker
func (worker *Worker) preparePosts() error {
	worker.log.Debug("Preparing Posts")
	//collect subreddits and posts
	subreddits := redditHelper.GetTopSubreddits()

	links, err := redditHelper.GetTopPostsForSubreddits(subreddits, 10)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare posts")
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
	return nil
}

//runRound handles a single round by parsing a post into the RoundMessage struct, sending it to all clients and spawning a handler for incoming answers
func (worker *Worker) runRound(round int) {
	worker.log.WithField("Round", round).Info("Starting Round")
	var wg sync.WaitGroup

	//get the post for this round
	post := worker.Posts[round]
	//if post is an html post use the HtmlContent, otherwise just the Content
	postText := post.Data.HtmlContent
	if postText == "" {
		postText = post.Data.Content
	}

	subreddits := preparePossibleAnswers(worker.Subreddits, post.Data.Subreddit)

	//Parse the post into our RoundMessage format and marshal it to json
	roundMessage := messages.NewRoundMessage()
	roundMessage.Payload.Number = round
	roundMessage.Payload.From = worker.RoundsTotal
	roundMessage.Payload.Post.Title = post.Data.Title
	roundMessage.Payload.Post.Content = postText
	roundMessage.Payload.Post.Type = post.GetType()
	roundMessage.Payload.Post.Url = post.Data.Url
	roundMessage.Payload.Subreddits = subreddits
	roundJson, err := json.Marshal(roundMessage)
	if err != nil {
		worker.log.WithError(err).
			Error("Unable to marshal round message to json")
		return
	}

	worker.Clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		if client.Terminated {
			return false
		}
		//flush messages still in channel before starting next round
		for i := 0; i < len(client.Message); i++ {
			worker.log.Trace(<-client.Message)
		}
		client.Send <- roundJson
		client.Blocked = false
		wg.Add(1)
		go worker.handleClientAnswer(client, post.Data.Subreddit, &wg)
		return true
	})
	wg.Wait()
}

func preparePossibleAnswers(possibleSubreddits []string, correctSubreddit string) []string {
	fakeSubreddits := make([]string, len(possibleSubreddits))
	copy(fakeSubreddits, possibleSubreddits)
	rand.Shuffle(len(fakeSubreddits), func(i, j int) {
		fakeSubreddits[i], fakeSubreddits[j] = fakeSubreddits[j], fakeSubreddits[i]
	})
	fakeSubreddits = utils.FilterString(fakeSubreddits, func(s string) bool {
		return s != correctSubreddit
	})
	subreddits := []string{correctSubreddit}
	subreddits = append(subreddits, fakeSubreddits[:3]...)
	rand.Shuffle(len(subreddits), func(i, j int) {
		subreddits[i], subreddits[j] = subreddits[j], subreddits[i]
	})
	return subreddits
}

func (worker *Worker) join(w http.ResponseWriter, r *http.Request, playerName string, playerUUID uuid.UUID) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		worker.log.Error("Unable to rejoin player")
	}
	gameClient := NewClient(conn, playerName, playerUUID, worker)
	worker.ClientCount++
	worker.log.WithField("Name", gameClient.Name).Debug("Player joined")
	worker.Clients.Store(gameClient.uuid.String(), gameClient)
	if len(worker.Posts) == 0 {
		err := worker.preparePosts()
		if err != nil {
			worker.State = Done
			return nil, errors.Wrap(err, "Failed to prepare posts. Game won't start!")
		}
	}
	return gameClient, nil
}

func (worker *Worker) sendScoreMessage() {
	msg := messages.NewScoreMessage()
	worker.Clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		score, _ := worker.Score.Load(client.uuid.String())
		msg.Payload.Scores[client.Name] = score.(int)
		return true
	})
	msgJson, err := json.Marshal(msg)
	if err != nil {
		worker.log.WithError(err).
			Error("Unable to marshal score message to json")
		return
	}
	worker.Clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		if !client.Terminated {
			client.Send <- msgJson
		}
		return true
	})
}

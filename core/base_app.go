package core

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/nanoteck137/kricketune/client/api"
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/player"
	"github.com/nanoteck137/kricketune/types"
)

var _ player.Queue = (*DwebbleQueue)(nil)

type List interface {
	GetName() string
	LoadTracks() ([]player.Track, error)
}

var _ List = (*Playlist)(nil)

type Playlist struct {
	client *api.Client

	Id   string
	Name string
}

func (p Playlist) GetName() string {
	return p.Name
}

func (p Playlist) LoadTracks() ([]player.Track, error) {
	items, err := p.client.GetMediaFromPlaylist(p.Id, api.GetMediaFromPlaylistBody{Shuffle: true}, api.Options{})
	if err != nil {
		return nil, err
	}

	tracks := make([]player.Track, len(items.Items))

	for i, t := range items.Items {
		tracks[i] = player.Track{
			Name:   t.Track.Name,
			Artist: t.Artists[0].Name,
			Album:  t.Album.Name,
			Uri:    t.MediaUrl,
		}
	}

	return tracks, nil
}

var _ List = (*Taglist)(nil)

type Taglist struct {
	client *api.Client

	Id   string
	Name string
}

func (p Taglist) GetName() string {
	return p.Name
}

func (p Taglist) LoadTracks() ([]player.Track, error) {
	items, err := p.client.GetMediaFromTaglist(p.Id, api.GetMediaFromTaglistBody{Shuffle: true}, api.Options{})
	if err != nil {
		return nil, err
	}

	tracks := make([]player.Track, len(items.Items))

	for i, t := range items.Items {
		tracks[i] = player.Track{
			Name:   t.Track.Name,
			Artist: t.Artists[0].Name,
			Album:  t.Album.Name,
			Uri:    t.MediaUrl,
		}
	}

	return tracks, nil
}

type DwebbleQueue struct {
	client *api.Client

	// TODO(patrik): Make private?
	Lists map[string]List

	mux    sync.RWMutex
	index  int
	tracks []player.Track
}

func NewDwebbleQueue(client *api.Client) *DwebbleQueue {
	return &DwebbleQueue{
		client: client,
		Lists:  map[string]List{},
	}
}

func GenerateCryptoID() string {
	bytes := make([]byte, 16)
	if _, err := cryptorand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

type QueueStatus struct {
	Index        int
	NumTracks    int
	CurrentTrack player.Track
}

func (q *DwebbleQueue) LoadList(list List) error {
	tracks, err := list.LoadTracks()
	if err != nil {
		return err
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	q.index = 0
	q.tracks = tracks

	return nil
}

func (q *DwebbleQueue) FetchLists() error {
	// q.mux.Lock()
	// defer q.mux.Unlock()
	clear(q.Lists)

	playlists, err := q.client.GetPlaylists(api.Options{})
	if err != nil {
		return err
	}

	for _, playlist := range playlists.Playlists {
		id := "playlist:" + playlist.Id
		name := fmt.Sprintf("Playlist - %s", playlist.Name)

		q.Lists[id] = Playlist{
			client: q.client,
			Id:     playlist.Id,
			Name:   name,
		}
	}

	taglists, err := q.client.GetTaglists(api.Options{})
	if err != nil {
		return err
	}

	for _, taglist := range taglists.Taglists {
		id := "taglist:" + taglist.Id
		name := fmt.Sprintf("Taglist - %s", taglist.Name)

		q.Lists[id] = Taglist{
			client: q.client,
			Id:     taglist.Id,
			Name:   name,
		}
	}

	return nil
}

func (q *DwebbleQueue) GetStatus() QueueStatus {
	q.mux.Lock()
	defer q.mux.Unlock()

	var currentTrack player.Track

	if len(q.tracks) > 0 {
		currentTrack = q.tracks[q.index]
	}

	return QueueStatus{
		Index:        q.index,
		NumTracks:    len(q.tracks),
		CurrentTrack: currentTrack,
	}
}

func (q *DwebbleQueue) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.index = 0
	q.tracks = nil
}

func (q *DwebbleQueue) CurrentTrack() (player.Track, bool) {
	q.mux.RLock()
	defer q.mux.RUnlock()

	if len(q.tracks) <= 0 {
		return player.Track{}, false
	}

	return q.tracks[q.index], true
}

func (q *DwebbleQueue) Next() {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.index++
	if q.index >= len(q.tracks) {
		q.index = 0
	}
}

func (q *DwebbleQueue) Prev() {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.index--
	if q.index < 0 {
		q.index = 0
	}
}

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	config *config.Config
	player *player.Player

	client *api.Client
	user   *User
	queue  *DwebbleQueue
}

func (app *BaseApp) User() *User {
	return nil
}

func (app *BaseApp) Queue() *DwebbleQueue {
	return app.queue
}

func (app *BaseApp) Player() *player.Player {
	return app.player
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	var err error

	app.player, err = player.New(app.config.AudioOutput)
	if err != nil {
		return fmt.Errorf("Failed to create audio player: %w", err)
	}

	app.queue.client.SetApiToken(app.config.ApiToken)

	user, err := app.queue.client.GetMe(api.Options{})
	if err == nil {
		app.user = &User{
			Username:        user.Username,
			DisplayName:     user.DisplayName,
			QuickPlaylistId: user.QuickPlaylist,
		}

		app.queue.FetchLists()
	}

	// NOTE(patrik): Setting default values
	app.player.SetVolume(1.0)
	app.player.SetMute(false)

	app.player.SetQueue(app.queue)

	// err = app.queue.LoadPlaylist(*user.QuickPlaylist)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	client := api.New(config.DwebbleAddress)
	queue := NewDwebbleQueue(client)

	return &BaseApp{
		config: config,
		queue:  queue,
	}
}

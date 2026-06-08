package core

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"sync"

	"github.com/nanoteck137/kricketune/client/api"
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/player"
	"github.com/nanoteck137/kricketune/tools/hook"
	"github.com/nanoteck137/kricketune/types"
)

type List interface {
	GetId() string
	GetName() string
	LoadTracks() ([]player.Track, error)
}

var _ List = (*Playlist)(nil)

type Playlist struct {
	client *api.Client

	Id     string
	FullId string
	Name   string

	Filter string
}

func (p *Playlist) GetId() string {
	return p.FullId
}

func (p *Playlist) GetName() string {
	return p.Name
}

func (p *Playlist) LoadTracks() ([]player.Track, error) {
	// items, err := p.client.GetMediaFromPlaylist(p.Id, api.GetMediaFromPlaylistBody{Shuffle: true}, api.Options{})
	items, err := p.client.GetPlaylistItems(p.Id, api.Options{
		Query: url.Values{
			"filter":  []string{p.Filter},
			"perPage": []string{"999999"},
		},
	})
	if err != nil {
		return nil, err
	}

	tracks := make([]player.Track, len(items.Items))

	for i, t := range items.Items {
		artists := make([]string, len(t.Artists))

		for i, artist := range t.Artists {
			artists[i] = artist.Name
		}

		// TODO(patrik): I don't handle the error here
		uri, _ := p.client.Url.StreamTrack(t.Id)

		tracks[i] = player.Track{
			Name:     t.Name,
			Artists:  artists,
			Album:    t.AlbumName,
			CoverUrl: t.CoverArt.Original,
			Uri:      uri.String(),
		}
	}

	return tracks, nil
}

var _ player.Queue = (*DwebbleQueue)(nil)

type DwebbleQueue struct {
	app    App
	client *api.Client

	// TODO(patrik): Make private?
	Lists         []List
	ListsById     map[string]List
	currentListId string

	mux    sync.RWMutex
	index  int
	tracks []player.Track
}

func (q *DwebbleQueue) addList(list List) {
	q.Lists = append(q.Lists, list)
	q.ListsById[list.GetId()] = list
}

func (q *DwebbleQueue) SetQueueIndex(index int) {
	if index < 0 {
		index = 0
	}

	if index >= len(q.tracks) {
		index = len(q.tracks) - 1
	}

	q.index = index
}

func NewDwebbleQueue(app App, client *api.Client) *DwebbleQueue {
	return &DwebbleQueue{
		app:       app,
		client:    client,
		Lists:     []List{},
		ListsById: map[string]List{},
		mux:       sync.RWMutex{},
		index:     0,
		tracks:    []player.Track{},
	}
}

type QueueStatus struct {
	Index        int
	NumTracks    int
	CurrentTrack player.Track

	CurrentListName string
	CurrentListId   string
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
	q.currentListId = list.GetId()

	q.app.OnQueueChanged().Call(context.TODO(), &OnQueueChangedEvent{})

	return nil
}

func (q *DwebbleQueue) FetchLists() error {
	q.mux.Lock()
	defer q.mux.Unlock()

	clear(q.Lists)

	filters, err := q.client.GetUserTrackFilters(q.app.User().Id, api.Options{})
	if err != nil {
		// TODO(patrik): Handle error
		return err
	}

	playlists, err := q.client.GetPlaylists(api.Options{})
	if err != nil {
		return err
	}

	for _, playlist := range playlists.Playlists {
		q.addList(&Playlist{
			client: q.client,
			Id:     playlist.Id,
			FullId: "playlist:" + playlist.Id,
			Name:   fmt.Sprintf("Playlist - %s", playlist.Name),
		})

		for _, filter := range filters.Filters {
			q.addList(&Playlist{
				client: q.client,
				Id:     playlist.Id,
				FullId: "playlist:" + playlist.Id + ":" + filter.FilterId,
				Name:   fmt.Sprintf("Playlist - %s (%s)", playlist.Name, filter.Name),
				Filter: filter.Filter,
			})
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

	var listName string
	if q.currentListId != "" {
		listName = q.ListsById[q.currentListId].GetName()
	}

	return QueueStatus{
		Index:           q.index,
		NumTracks:       len(q.tracks),
		CurrentTrack:    currentTrack,
		CurrentListName: listName,
		CurrentListId:   q.currentListId,
	}
}

func (q *DwebbleQueue) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.index = 0
	q.tracks = nil
}

func (q *DwebbleQueue) ClearQueue() {
	q.Clear()
	q.app.OnQueueChanged().Call(context.TODO(), &OnQueueChangedEvent{})
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

func (q *DwebbleQueue) CloneTracks() []player.Track {
	q.mux.Lock()
	defer q.mux.Unlock()

	return slices.Clone(q.tracks)
}

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	config *config.Config
	player *player.Player

	client *api.Client
	user   *api.GetMe
	queue  *DwebbleQueue

	onQueueChanged *hook.Hook[*OnQueueChangedEvent]
}

func NewBaseApp(config *config.Config) *BaseApp {
	client := api.New(config.ApiAddress)

	app := &BaseApp{
		config: config,
	}

	queue := NewDwebbleQueue(app, client)
	app.queue = queue

	app.onQueueChanged = hook.NewHook[*OnQueueChangedEvent]("onQueueChanged")

	return app
}

func (app *BaseApp) Bootstrap() error {
	var err error

	app.player, err = player.New(app.config.AudioOutput)
	if err != nil {
		return fmt.Errorf("Failed to create audio player: %w", err)
	}

	// Setup the API Token
	app.queue.client.Headers.Set("X-Api-Token", app.config.ApiToken)

	user, err := app.queue.client.GetMe(api.Options{})
	if err == nil {
		app.user = user

		err = app.queue.FetchLists()
		if err != nil {
			return err
		}
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

func (app *BaseApp) User() *api.GetMe {
	return app.user
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

func (app *BaseApp) OnQueueChanged() *hook.Hook[*OnQueueChangedEvent] {
	return app.onQueueChanged
}

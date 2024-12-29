package core

import (
	"fmt"
	"sync"

	"github.com/nanoteck137/kricketune/client/api"
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/player"
	"github.com/nanoteck137/kricketune/types"
)

var _ player.Queue = (*DwebbleQueue)(nil)

type DwebbleQueue struct {
	client *api.Client

	mux    sync.RWMutex
	index  int
	tracks []player.Track
}

func NewDwebbleQueue(client *api.Client) *DwebbleQueue {
	return &DwebbleQueue{
		client: client,
	}
}

func (q *DwebbleQueue) LoadFilter(filter, sort string) error {
	q.mux.Lock()
	defer q.mux.Unlock()

	tracks, err := q.client.GetTracks(api.Options{
		QueryParams: map[string]string{
			"filter": filter,
			"sort":   sort,
			"perPage": "500",
		},
	})
	if err != nil {
		return err
	}

	for _, t := range tracks.Tracks {
		q.tracks = append(q.tracks, player.Track{
			Name:   t.Name.Default,
			Artist: t.ArtistName.Default,
			Album:  t.AlbumName.Default,
			Uri:    t.MobileMediaUrl,
		})
	}

	return nil
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
	queue  *DwebbleQueue
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

	// NOTE(patrik): Setting default values
	app.player.SetVolume(1.0)
	app.player.SetMute(false)

	app.player.SetQueue(app.queue)

	if len(app.Config().FilterSets) > 0 {
		set := app.Config().FilterSets[0]

		err := app.queue.LoadFilter(set.Filter, set.Sort)
		if err != nil {
			return err
		}
	}

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

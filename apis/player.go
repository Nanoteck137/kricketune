package apis

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/kricketune/player"
	"github.com/nanoteck137/kricketune/tools/broker"
	"github.com/nanoteck137/kricketune/utils"
	"github.com/nanoteck137/pyrin"
)

type Set struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

type Sets struct {
	Sets []Set `json:"sets"`
}

type Track struct {
	Name     string   `json:"name"`
	Artists  []string `json:"artists"`
	Album    string   `json:"album"`
	CoverUrl string   `json:"coverUrl"`
}

func ConvertPlayerTrackToTrack(t player.Track) Track {
	return Track{
		Name:     t.Name,
		Artists:  utils.FixNilArrayToEmpty(t.Artists),
		Album:    t.Album,
		CoverUrl: t.CoverUrl,
	}
}

type Status struct {
	CurrentTrack Track `json:"currentTrack"`

	IsPlaying bool    `json:"isPlaying"`
	Volume    float32 `json:"volume"`
	Mute      bool    `json:"mute"`

	QueueIndex int `json:"queueIndex"`
	NumTracks  int `json:"numTracks"`

	Position int64 `json:"position"`
	Duration int64 `json:"duration"`
}

type List struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetLists struct {
	Lists []List `json:"lists"`
}

type SeekBody struct {
	Skip int `json:"skip"`
}

func InstallPlayerHandlers(app core.App, group pyrin.Group) {
	b := broker.NewBroker()
	b.Start()

	// TODO(patrik): Use http.Method*
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetStatus",
			Method:       "GET",
			Path:         "/player/status",
			ResponseType: Status{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				queueStatus := app.Queue().GetStatus()
				position, duration := app.Player().GetPosition()

				track := ConvertPlayerTrackToTrack(queueStatus.CurrentTrack)

				return Status{
					CurrentTrack: track,
					IsPlaying:    app.Player().IsPlaying(),
					Volume:       app.Player().GetVolume(),
					Mute:         app.Player().GetMute(),
					QueueIndex:   queueStatus.Index,
					NumTracks:    queueStatus.NumTracks,
					Position:     position / int64(time.Millisecond),
					Duration:     duration / int64(time.Millisecond),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetLists",
			Method:       "GET",
			Path:         "/player/lists",
			ResponseType: GetLists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				queue := app.Queue()

				res := GetLists{
					Lists: make([]List, 0, len(queue.Lists)),
				}

				for id, list := range queue.Lists {
					res.Lists = append(res.Lists, List{
						Id:   id,
						Name: list.GetName(),
					})
				}

				sort.Slice(res.Lists, func(i, j int) bool {
					return res.Lists[i].Name < res.Lists[j].Name
				})

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "LoadList",
			Method: http.MethodPost,
			Path:   "/player/lists/:id",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				queue := app.Queue()

				list, exists := queue.Lists[id]
				if !exists {
					// TODO(patrik): Error
					return nil, errors.New("No list with id")
				}

				err := queue.LoadList(list)
				if err != nil {
					return nil, err
				}

				app.Player().Start()

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Play",
			Method: "POST",
			Path:   "/player/play",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().Play()
				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Pause",
			Method: "POST",
			Path:   "/player/pause",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().Pause()
				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Next",
			Method: "POST",
			Path:   "/player/next",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().PrepareChange()
				app.Player().NextTrack()
				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Prev",
			Method: "POST",
			Path:   "/player/prev",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().PrepareChange()
				app.Player().PrevTrack()
				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RewindTrack",
			Method: "POST",
			Path:   "/player/rewindTrack",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().RewindTrack()

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "Seek",
			Method:   "POST",
			Path:     "/player/seek",
			BodyType: SeekBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[SeekBody](c)
				if err != nil {
					return nil, err
				}

				app.Player().Seek(time.Duration(body.Skip) * time.Second)

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearQueue",
			Method: "POST",
			Path:   "/player/clearQueue",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// app.Player().ClearQueue()
				return nil, nil
			},
		},

		pyrin.NormalHandler{
			Name:        "SseHandler",
			Method:      http.MethodGet,
			Path:        "/player/sse",
			HandlerFunc: func(c pyrin.Context) error {
				b.ServeHTTP(c.Response(), c.Request())
				return nil
			},
		},
	)
}

package apis

import (
	"net/http"
	"strconv"

	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/pyrin"
)

type Set struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

type Sets struct {
	Sets []Set `json:"sets"`
}

type Status struct {
	TrackName   string `json:"trackName"`
	TrackArtist string `json:"trackArtist"`
	TrackAlbum  string `json:"trackAlbum"`

	IsPlaying bool    `json:"isPlaying"`
	Volume    float32 `json:"volume"`
	Mute      bool    `json:"mute"`

	QueueIndex int `json:"queueIndex"`
	NumTracks  int `json:"numTracks"`
}

func InstallPlayerHandlers(app core.App, group pyrin.Group) {
	// TODO(patrik): Use http.Method*
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetSets",
			Method:   http.MethodGet,
			Path:     "/player/sets",
			DataType: Sets{},
			HandlerFunc: func(c pyrin.Context) (any, error) {

				res := Sets{
					Sets: make([]Set, len(app.Config().FilterSets)),
				}

				for i, set := range app.Config().FilterSets {
					res.Sets[i] = Set{
						Name:  set.Name,
						Index: i,
					}
				}

				return res, nil
			},
		},
		pyrin.ApiHandler{
			Name:     "ChangeSet",
			Method:   http.MethodPost,
			Path:     "/player/sets/:index",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				indexStr := c.Param("index")
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Add index checks
				set := app.Config().FilterSets[index]

				app.Queue().Clear()
				app.Queue().LoadFilter(set.Filter, set.Sort)
				app.Player().Start()

				return nil, nil
			},
		},
		pyrin.ApiHandler{
			Name:     "GetStatus",
			Method:   "GET",
			Path:     "/player/status",
			DataType: Status{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				currentTrack, _ := app.Queue().CurrentTrack()

				res := Status{
					TrackName:   currentTrack.Name,
					TrackArtist: currentTrack.Artist,
					TrackAlbum:  currentTrack.Album,
					IsPlaying:   app.Player().IsPlaying(),
					Volume:      app.Player().GetVolume(),
					Mute:        app.Player().GetMute(),
					// QueueIndex:  app.Player().CurrentQueueIndex(),
					// NumTracks:   app.Player().NumTracks(),
				}

				return res, nil
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
			Name:   "ClearQueue",
			Method: "POST",
			Path:   "/player/clearQueue",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// app.Player().ClearQueue()
				return nil, nil
			},
		},
	)
}

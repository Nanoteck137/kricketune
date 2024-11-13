package apis

import (
	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/pyrin"
)

func InstallPlayerHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:        "Play",
			Method:      "POST",
			Path:        "/player/play",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().Play()
				return nil, nil
			},
		},
		pyrin.ApiHandler{
			Name:        "Pause",
			Method:      "POST",
			Path:        "/player/pause",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().Pause()
				return nil, nil
			},
		},
		pyrin.ApiHandler{
			Name:        "Next",
			Method:      "POST",
			Path:        "/player/next",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().PrepareChange()
				app.Player().NextTrack()
				return nil, nil
			},
		},
		pyrin.ApiHandler{
			Name:        "Prev",
			Method:      "POST",
			Path:        "/player/prev",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().PrepareChange()
				app.Player().PrevTrack()
				return nil, nil
			},
		},
		pyrin.ApiHandler{
			Name:        "ClearQueue",
			Method:      "POST",
			Path:        "/player/clearQueue",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				app.Player().ClearQueue()
				return nil, nil
			},
		},
	)
}

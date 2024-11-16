package apis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nanoteck137/kricketune"
	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/pyrin"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)

	g.Register(pyrin.NormalHandler{
		Method:      http.MethodGet,
		Path:        "/sse",
		HandlerFunc: func(c pyrin.Context) error {
			w := c.Response()
			r := c.Request()

			// Set http headers required for SSE
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			// You may need this locally for CORS requests
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Create a channel for client disconnection
			clientGone := r.Context().Done()

			rc := http.NewResponseController(w)
			t := time.NewTicker(time.Second)
			defer t.Stop()
			for {
				select {
				case <-clientGone:
					fmt.Println("Client disconnected")
					return nil
				case <-t.C:
					// Send an event to the client
					// Here we send only the "data" field, but there are few others
					_, err := fmt.Fprintf(w, "data: The time is %s\n\n", time.Now().Format(time.UnixDate))
					if err != nil {
						return nil
					}
					err = rc.Flush()
					if err != nil {
						return nil
					}
				}
			}
		},
	})
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: kricketune.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	return s, nil
}

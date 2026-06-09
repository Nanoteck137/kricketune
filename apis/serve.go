package apis

import (
	"net/http"
	"os"

	"github.com/nanoteck137/kricketune"
	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/pyrin"
)

func RegisterStaticHandlers(app core.App, g pyrin.Group) {
	g.Register(
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/static/*",
			HandlerFunc: func(c pyrin.Context) error {
				// TODO(patrik): Fix this
				f := os.DirFS("./render/static")
				fs := http.StripPrefix("/static", http.FileServerFS(f))

				fs.ServeHTTP(c.Response(), c.Request())

				return nil
			},
		},
	)

	// TODO(patrik): I don't like this
	if app != nil {
		webDir := app.Config().WebDir
		if webDir != "" {
			g.Register(
				// TODO(patrik): Fix this
				pyrin.SpaHandler(os.DirFS(webDir), "index.html"),
			)
		}
	}
}

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)
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

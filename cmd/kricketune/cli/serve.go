package cli

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/nanoteck137/kricketune/apis"
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/kricketune/core/log"
	"github.com/nanoteck137/kricketune/player"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {

		gst.Init(nil)

		p, err := player.New()
		if err != nil {
			log.Fatal("Failed to create player", "err", err)
		}

		app := core.NewBaseApp(&config.LoadedConfig, p)

		p.AddTrack(player.Track{
			Name:   "Test",
			Artist: "Test",
			Uri:    "https://dwebble.nanoteck137.net/files/tracks/original/infinity-on-high-oykr7v5k/02-the-take-over-the-breaks-over-1726248927.opus",
		})

		p.AddTrack(player.Track{
			Name:   "Test",
			Artist: "Test",
			Uri:    "https://dwebble.nanoteck137.net/files/tracks/original/115-feat-apricot-cae4bw4p/1-115-feat-apricot-1726248749.opus",
		})

		p.PrepareChange()
		p.Reset()

		err = player.Launch(p)
		if err != nil {
			log.Fatal("Failed to launch player", "err", err)
		}

		err = app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		e, err := apis.Server(app)
		if err != nil {
			log.Fatal("Failed to create server", "err", err)
		}

		err = e.Start(app.Config().ListenAddr)
		if err != nil {
			log.Fatal("Failed to start server", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

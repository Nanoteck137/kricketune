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

		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		err = player.Launch(app.Player())
		if err != nil {
			log.Fatal("Failed to launch player", "err", err)
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

package cli

import (
	"github.com/nanoteck137/kricketune"
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/core/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     kricketune.AppName,
	Version: kricketune.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(kricketune.VersionTemplate(kricketune.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}

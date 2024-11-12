package config

import (
	"os"

	"github.com/nanoteck137/kricketune"
	"github.com/nanoteck137/kricketune/core/log"
	"github.com/nanoteck137/kricketune/types"
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr      string `mapstructure:"listen_addr"`
	DataDir         string `mapstructure:"data_dir"`
}

func (c *Config) WorkDir() types.WorkDir {
	return types.WorkDir(c.DataDir)
}

func setDefaults() {
	viper.SetDefault("listen_addr", ":3000")
	viper.BindEnv("data_dir")
}

func validateConfig(config *Config) {
	hasError := false

	validate := func(expr bool, msg string) {
		if expr {
			log.Error("Config Validation", "err", msg)
			hasError = true
		}
	}

	// NOTE(patrik): Has default value, here for completeness
	validate(config.ListenAddr == "", "listen_addr needs to be set")
	validate(config.DataDir == "", "data_dir needs to be set")

	if hasError {
		log.Fatal("Config not valid")
	}
}

var ConfigFile string
var LoadedConfig Config

func InitConfig() {
	setDefaults()

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(kricketune.AppName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("Failed to load config", "err", err)
	}

	err = viper.Unmarshal(&LoadedConfig)
	if err != nil {
		log.Error("Failed to unmarshal config: ", err)
		os.Exit(-1)
	}

	log.Debug("Current Config", "config", LoadedConfig)

	validateConfig(&LoadedConfig)
}

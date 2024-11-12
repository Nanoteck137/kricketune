package core

import (
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	config   *config.Config
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}

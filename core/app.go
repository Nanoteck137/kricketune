package core

import (
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	Config() *config.Config

	WorkDir() types.WorkDir

	Bootstrap() error
}
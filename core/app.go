package core

import (
	"github.com/nanoteck137/kricketune/config"
	"github.com/nanoteck137/kricketune/player"
	"github.com/nanoteck137/kricketune/types"
)

type User struct {
	Username        string
	DisplayName     string
	QuickPlaylistId *string
}

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	Config() *config.Config
	Player() *player.Player
	User() *User
	Queue() *DwebbleQueue

	WorkDir() types.WorkDir

	Bootstrap() error
}

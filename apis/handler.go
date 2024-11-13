package apis

import (
	"github.com/nanoteck137/kricketune/core"
	"github.com/nanoteck137/pyrin"
)

func InstallHandlers(app core.App, g pyrin.Group) {
	InstallPlayerHandlers(app, g)
}

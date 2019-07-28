// common.go
package routes

import (
	"github.com/tomicooler/zukebox/models"
	"github.com/tomicooler/zukebox/zuke"
)

var Cache zuke.CacheManager
var Tracks models.TracksManager
var Control models.Control
var PlayerCtx zuke.PlayerContext

func init() {
	Tracks = models.NewTracksManager()
	Cache = zuke.NewCacheManager()
	Control = models.NewControl()
	PlayerCtx = zuke.NewPlayerContext(&Tracks, &Cache, &Control)
	go PlayerCtx.Start()
}

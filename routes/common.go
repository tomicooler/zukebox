// common.go
package routes

import (
	"github.com/tomicooler/zukebox/models"
	"github.com/tomicooler/zukebox/zuke"
)

var Tracks models.TracksManager
var Control models.Control
var PlayerCtx zuke.PlayerContext

func init() {
	Tracks = models.NewTracksManager()
	Control = models.NewControl()
	PlayerCtx = zuke.NewPlayerContext(&Tracks, &Control, &Ydl)
	go PlayerCtx.Start()
}

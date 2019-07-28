// player.go
package zuke

import (
	"log"
	"math"
	"sync"
	"time"

	vlc "github.com/adrg/libvlc-go"
	"github.com/tomicooler/zukebox/models"
)

type PlayerContext struct {
	tracks  *models.TracksManager
	cache   *CacheManager
	control *models.Control
	player  *vlc.Player
	media   *vlc.Media
	manager *vlc.EventManager
	ticker  *time.Ticker
	track   models.Track
	mutex   sync.Mutex
}

func NewPlayerContext(tracks *models.TracksManager, cache *CacheManager, control *models.Control) PlayerContext {
	if err := vlc.Init("--no-video"); err != nil {
		log.Fatal(err)
	}
	ctx := PlayerContext{tracks: tracks, cache: cache, control: control, player: nil, media: nil, manager: nil, ticker: nil, track: models.Track{}, mutex: sync.Mutex{}}

	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	ctx.player = player

	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}
	ctx.manager = manager

	ctx.ticker = time.NewTicker(250 * time.Millisecond)

	return ctx
}

func (ctx *PlayerContext) Release() {
	if ctx.player != nil {
		ctx.player.Stop()
		ctx.player.Release()
	}
	if ctx.media != nil {
		ctx.media.Release()
	}

	vlc.Release()
}

func (ctx *PlayerContext) Start() {
	for range ctx.ticker.C {
		ctx.UpdateControl()
		ctx.Next()
	}
}

func (ctx *PlayerContext) Seek(time float64) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if ctx.player != nil {
		log.Println("Seek", time, ctx.track.Duration)
		ctx.player.SetMediaPosition(float32(math.Max(0, math.Min(time, ctx.track.Duration)) / ctx.track.Duration))
	}
}

func (ctx *PlayerContext) SetVolume(volume float64) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if ctx.player != nil {
		log.Println("SetVolume", volume)
		ctx.player.SetVolume(int(math.Max(0, math.Min(volume, 100))))
	}
}

func (ctx *PlayerContext) Next() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	if ctx.player == nil {
		return
	}

	state, err := ctx.player.MediaState()
	if err != nil {
		return
	}

	if !ctx.control.Enabled {
		if state == vlc.MediaPlaying {
			log.Println("Pause")
			ctx.player.SetPause(true)
		}
		return
	}

	if state == vlc.MediaPlaying {
		return
	}

	if state == vlc.MediaEnded || state == vlc.MediaError || state == vlc.MediaNothingSpecial {
		track, err := ctx.tracks.Pop()
		if err != nil {
			return
		}
		ctx.track = track

		log.Println("Load new track", ctx.track)

		if ctx.media != nil {
			ctx.media.Release()
		}

		if len(track.Message) > 0 {
			lang := track.Lang
			if len(lang) == 0 {
				lang = "en"
			}

			url := GetUrlFromSpeech(track.Message, lang)
			media, err := ctx.player.LoadMediaFromURL(url)
			if err != nil {
				log.Println("Could not load media from url", err, url)
				return
			}
			defer media.Release()

			err = ctx.player.Play()
			if err != nil {
				log.Println("Could not play media from url", err)
				return
			}

			quit := make(chan struct{})
			waitForFinish := func(event vlc.Event, userData interface{}) {
				close(quit)
			}

			eventID, err := ctx.manager.Attach(vlc.MediaPlayerEndReached, waitForFinish, nil)
			if err != nil {
				log.Println("Could not attach event for MediaPlayerEndReached", err)
				return
			}
			defer ctx.manager.Detach(eventID)

			<-quit
		}

		media, err := ctx.player.LoadMediaFromPath(ctx.cache.TrackPath(track.ID))
		if err != nil {
			log.Println("Could not load media from path", err)
			return
		}
		ctx.media = media

		err = ctx.player.Play()
		if err != nil {
			log.Println("Could not play media from path", err)
			return
		}
	} else {
		if state == vlc.MediaPaused {
			log.Println("Unpause")
			ctx.player.SetPause(false)
		}
	}
}

func (ctx *PlayerContext) UpdateControl() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	if ctx.player == nil {
		return
	}

	if state, err := ctx.player.MediaState(); err == nil {
		if state == vlc.MediaPlaying {
			ctx.control.SetPlaying(true)
			if pos, err := ctx.player.MediaPosition(); err == nil {
				ctx.control.SetTime(float64(int(float64(pos) * ctx.track.Duration)))
			}
		} else if state == vlc.MediaPaused {
			ctx.control.SetPlaying(false)
		} else if state == vlc.MediaEnded || state == vlc.MediaError {
			if ctx.track.ID != "" {
				ctx.tracks.AddRecentTrack(ctx.track)
			}
			ctx.track = models.Track{}
			ctx.control.SetTime(0)
		}

		if vol, err := ctx.player.Volume(); err == nil {
			ctx.control.SetVolume(float64(vol))
		}
	}

	ctx.control.SetTrack(ctx.track)
}

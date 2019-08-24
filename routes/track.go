// track.go
package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tomicooler/zukebox/models"
	"github.com/tomicooler/zukebox/zuke"
)

var ydl zuke.YoutubeDL
var ticker *time.Ticker

func init() {
	ydl = zuke.NewYoutubeDl()
	ticker = time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			ydl.Update()
		}
	}()
}

func GetATrack(w http.ResponseWriter, r *http.Request) {
	trackID, err := strconv.Atoi(chi.URLParam(r, "trackID"))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	track, err := Tracks.GetTrack(trackID)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.JSON(w, r, track)
}

func DeleteTrack(w http.ResponseWriter, r *http.Request) {
	trackID, err := strconv.Atoi(chi.URLParam(r, "trackID"))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := Tracks.DeleteTrack(trackID); err != nil {
		render.Render(w, r, ErrNotFound)
	}
}

func CreateTrack(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	track := models.Track{}
	if err := decoder.Decode(&track); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if cachedTrack, err := Cache.GetTrack(track.Url); err == nil {
		cachedTrack.User = track.User
		cachedTrack.Message = track.Message
		cachedTrack.Lang = track.Lang
		track = cachedTrack
	} else {
		info, err := ydl.GetTrackInfo(track.Url)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if info.IsLive {
			render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: "Live track is not supported!"})
			return
		}

		track.ID = info.ID
		track.Title = info.Title
		track.Thumbnail = info.Thumbnail
		track.Duration = info.Duration

		Cache.RemoveOldies()

		if err := ydl.DownloadTrack(track.Url, Cache.CachePath); err != nil {
			render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: "Could not download track!"})
			return
		}

		if err := Cache.StoreTrack(track); err != nil {
			render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: "Could not store track!"})
			return
		}
	}

	if err := Tracks.AddTrack(track); err != nil {
		render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: err.Error()})
		return
	}

	render.JSON(w, r, track)
}

func GetAllTrack(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Tracks.GetTracks())
}

func GetARecentTrack(w http.ResponseWriter, r *http.Request) {
	trackID, err := strconv.Atoi(chi.URLParam(r, "trackID"))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	track, err := Tracks.GetRecentTrack(trackID)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.JSON(w, r, track)
}

func GetAllRecentTrack(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Tracks.GetRecentTracks())
}

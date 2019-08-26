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

var Ydl zuke.YoutubeDL
var ticker *time.Ticker

func init() {
	Ydl = zuke.NewYoutubeDl()
	ticker = time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			Ydl.Update()
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

	info, err := Ydl.GetTrackInfo(track.Url)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if info.PlayerResponse.VideoDetails.IsLive {
		render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: "Live track is not supported!"})
		return
	}

	track.ID = info.ID
	track.Title = info.Title
	if len(info.PlayerResponse.VideoDetails.Thumbnail.Thumbnails) > 0 {
		track.Thumbnail = info.PlayerResponse.VideoDetails.Thumbnail.Thumbnails[0].Url
	}
	duration, err := strconv.ParseFloat(info.Duration, 64)
	if err != nil {
		render.Render(w, r, &ErrResponse{HTTPStatusCode: 403, StatusText: "Could not parse duration!"})
		return
	}

	track.Duration = duration

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

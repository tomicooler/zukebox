// track.go
package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func GetControl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Control.Get())
}

func PatchControl(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var dict map[string]interface{}
	decoder.Decode(&dict)

	if p, ok := dict["playing"]; ok {
		if play, ok := p.(bool); ok {
			Control.SetEnabled(play)
		} else {
			render.Render(w, r, ErrInvalidRequest(errors.New("playing must be a boolean value")))
			return
		}
	}

	if t, ok := dict["time"]; ok {
		if time, ok := t.(float64); ok {
			PlayerCtx.Seek(time)
		} else {
			render.Render(w, r, ErrInvalidRequest(errors.New("time must be a float value")))
			return
		}
	}

	if v, ok := dict["volume"]; ok {
		if volume, ok := v.(float64); ok {
			PlayerCtx.SetVolume(volume)
		} else {
			render.Render(w, r, ErrInvalidRequest(errors.New("volume must be a float value")))
			return
		}
	}

	render.JSON(w, r, Control.Get())
}

// control.go
package models

import (
	"math"
	"sync"
)

type Control struct {
	Playing bool       `json:"playing,omitempt"`
	Volume  float64    `json:"volume,omitempt"`
	Time    float64    `json:"time,omitempt"`
	Track   Track      `json:"track,omitempt"`
	Enabled bool       `json:"-"`
	mutex   sync.Mutex `json:"-"`
}

func NewControl() Control {
	return Control{Playing: false, Volume: 100, Time: 0, Track: Track{}, Enabled: true, mutex: sync.Mutex{}}
}

func (ctrl *Control) Get() Control {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	return *ctrl
}

func (ctrl *Control) SetPlaying(playing bool) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.Playing = playing
}

func (ctrl *Control) SetTrack(track Track) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.Track = track
}

func (ctrl *Control) SetEnabled(enabled bool) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.Enabled = enabled
}

func (ctrl *Control) SetTime(time float64) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.Time = time
}

func (ctrl *Control) SetVolume(volume float64) {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.Volume = math.Max(0, math.Min(volume, 100))
}

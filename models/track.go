// models.go
package models

import (
	"errors"
	"sync"
)

type Track struct {
	ID        string  `json:id`
	Url       string  `json:"url"`
	User      string  `json:"user"`
	Message   string  `json:"message,omitempt"`
	Lang      string  `json:"lang,omitempt"`
	Duration  float64 `json:"duration,omitempt"`
	Title     string  `json:"title,omitempt"`
	Thumbnail string  `json:"thumbnail,omitempt"`
}

type Tracks map[string][]Track

type TracksManager struct {
	tracks       []Track
	recentTracks []Track
	maxSize      int
	mutex        sync.Mutex
}

func NewTracksManager() TracksManager {
	return TracksManager{tracks: make([]Track, 0), recentTracks: make([]Track, 0), maxSize: 100, mutex: sync.Mutex{}}
}

func (tracks *TracksManager) Pop() (Track, error) {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	if len(tracks.tracks) == 0 {
		return Track{}, errors.New("Empty")
	}

	track := tracks.tracks[0]
	tracks.tracks = append(tracks.tracks[:0], tracks.tracks[1:]...)

	return track, nil
}

func (tracks *TracksManager) GetTrack(index int) (Track, error) {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	if index >= len(tracks.tracks) {
		return Track{}, errors.New("Not found!")
	}

	return tracks.tracks[index], nil
}

func (tracks *TracksManager) GetTracks() Tracks {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()
	return Tracks{"tracks": tracks.tracks}
}

func (tracks *TracksManager) AddTrack(track Track) error {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	if len(tracks.tracks) >= tracks.maxSize {
		return errors.New("Max tracks reached!")
	}

	tracks.tracks = append(tracks.tracks, track)

	return nil
}

func (tracks *TracksManager) DeleteTrack(index int) error {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	if index >= len(tracks.tracks) {
		return errors.New("Not found!")
	}

	tracks.tracks = append(tracks.tracks[:index], tracks.tracks[index+1:]...)
	return nil
}

func (tracks *TracksManager) GetRecentTrack(index int) (Track, error) {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	if index >= len(tracks.recentTracks) {
		return Track{}, errors.New("Not found!")
	}

	return tracks.tracks[index], nil
}

func (tracks *TracksManager) GetRecentTracks() Tracks {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()
	return Tracks{"tracks": tracks.recentTracks}
}

func (tracks *TracksManager) AddRecentTrack(track Track) {
	tracks.mutex.Lock()
	defer tracks.mutex.Unlock()

	tracks.recentTracks = append([]Track{track}, tracks.recentTracks...)

	if len(tracks.recentTracks) > tracks.maxSize {
		tracks.recentTracks = tracks.recentTracks[:len(tracks.recentTracks)-1]
	}
}

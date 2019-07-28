// youtube.go
package zuke

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
)

type TrackInfo struct {
	ID        string  `json:"id"`
	Thumbnail string  `json:"thumbnail"`
	Duration  float64 `json:"duration"`
	Title     string  `json:"title"`
	IsLive    bool    `json:"is_live"`
}

type YoutubeDL struct {
	YoutubeDLPath string
}

func NewYoutubeDl() YoutubeDL {
	ydl := path.Join(binPath(), "ydl")
	os.MkdirAll(path.Dir(ydl), 0755)
	if _, err := os.Stat(ydl); os.IsNotExist(err) {
		DownloadFile(ydl, "https://yt-dl.org/downloads/latest/youtube-dl")
		os.Chmod(ydl, 0755)
	} else {
		updater := YoutubeDL{YoutubeDLPath: ydl}
		updater.Update()
	}

	return YoutubeDL{YoutubeDLPath: ydl}
}

func (ydl *YoutubeDL) GetTrackInfo(url string) (TrackInfo, error) {
	output, err := ydl.run([]string{"--dump-json", "--no-playlist", url})
	if err != nil {
		return TrackInfo{}, err
	}

	info := TrackInfo{}
	if err := json.Unmarshal(output, &info); err != nil {
		return TrackInfo{}, err
	}

	return info, nil
}

func (ydl *YoutubeDL) DownloadTrack(url string, outputdir string) error {
	if _, err := ydl.run([]string{"--extract-audio", "--no-playlist", "--output", path.Join(outputdir, "%(id)s.%(ext)s"), "--audio-format", "opus", url}); err != nil {
		return err
	}

	return nil
}

func (ydl *YoutubeDL) Update() error {
	if _, err := ydl.run([]string{"--update"}); err != nil {
		return err
	}

	return nil
}

func (ydl *YoutubeDL) run(args []string) ([]byte, error) {
	cmd := exec.Command(ydl.YoutubeDLPath, args...)
	return cmd.Output()
}

func binPath() string {
	if cachePath, err := os.UserCacheDir(); err == nil {
		return path.Join(cachePath, "zukebox", "bin")
	}
	return path.Join("tmp", "zukebox", "bin")
}

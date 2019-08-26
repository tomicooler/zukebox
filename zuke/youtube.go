// youtube.go
package zuke

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Thumbnail struct {
	Url string `json:"url"`
}

type Thumbnails struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
}

type VideoDetails struct {
	Thumbnail Thumbnails `json:"thumbnail"`
	IsLive    bool       `json:"isLiveContent"`
}

type PlayerResponse struct {
	VideoDetails VideoDetails `json:"videoDetails"`
}

type TrackInfo struct {
	ID             string         `json:"video_id"`
	Title          string         `json:"title"`
	Duration       string         `json:"length_seconds"`
	PlayerResponse PlayerResponse `json:"player_response"`
}

type YoutubeDL struct {
	YoutubeDLPath string
}

func NewYoutubeDl() YoutubeDL {
	npmPath := npmPath()
	ydl := path.Join(npmPath, "bin", "ytdl")
	os.MkdirAll(npmPath, 0755)
	if _, err := os.Stat(ydl); os.IsNotExist(err) {
		cmd := exec.Command("npm", []string{"install", "ytdl", "--prefix", npmPath, "-g"}...)
		if out, err := cmd.Output(); err != nil {
			log.Fatal("Could not install ytdl npm package", err, out)
		}
	} else {
		updater := YoutubeDL{YoutubeDLPath: ydl}
		updater.Update()
	}

	return YoutubeDL{YoutubeDLPath: ydl}
}

func (ydl *YoutubeDL) GetTrackInfo(url string) (TrackInfo, error) {
	output, err := ydl.run([]string{"--info-json", url})
	if err != nil {
		fmt.Println("error ", output, err)
		return TrackInfo{}, err
	}

	info := TrackInfo{}
	if err := json.Unmarshal(output, &info); err != nil {
		return TrackInfo{}, err
	}

	return info, nil
}

func (ydl *YoutubeDL) GetTrackUrl(url string) (string, error) {
	out, err := ydl.run([]string{"--filter", "audioonly", url, "--print-url"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), err
}

func (ydl *YoutubeDL) Update() error {
	cmd := exec.Command("npm", []string{"update", "ytdl", "--prefix", npmPath(), "-g"}...)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func (ydl *YoutubeDL) run(args []string) ([]byte, error) {
	cmd := exec.Command(ydl.YoutubeDLPath, args...)
	return cmd.Output()
}

func npmPath() string {
	if cachePath, err := os.UserCacheDir(); err == nil {
		return path.Join(cachePath, "zukebox", "npm")
	}
	return path.Join("tmp", "zukebox", "npm")
}

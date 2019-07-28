// cache.go
package zuke

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tomicooler/zukebox/models"
)

type CacheManager struct {
	CachePath string
	MaxSize   int64
	ExtInfo   string
	ExtTrack  string
	IdRegexp  *regexp.Regexp
}

func NewCacheManager() CacheManager {
	cPath := cachePath()
	os.MkdirAll(cPath, 0755)
	r, _ := regexp.Compile(`^(?:http(?:s)?://)?(?:www\.)?(?:m\.)?(?:youtu\.be/|youtube\.com/(?:(?:watch)?\?(?:.*&)?v(?:i)?=|(?:embed|v|vi|user)/))([^\?&\"'>]+)`)
	return CacheManager{CachePath: cPath, MaxSize: 1000, ExtInfo: ".json", ExtTrack: ".opus", IdRegexp: r}
}

func (cache *CacheManager) GetID(url string) (string, error) {
	groups := cache.IdRegexp.FindStringSubmatch(url)
	if len(groups) != 2 {
		return "", errors.New("No ID found")
	}
	return groups[1], nil
}

func (cache *CacheManager) GetTrack(url string) (models.Track, error) {
	id, err := cache.GetID(url)
	if err != nil {
		return models.Track{}, err
	}

	if !cache.IsCached(id) {
		return models.Track{}, errors.New("Not cached yet")
	}

	reader, err := os.Open(cache.InfoPath(id))
	if err != nil {
		return models.Track{}, err
	}

	decoder := json.NewDecoder(reader)
	track := models.Track{}
	if err := decoder.Decode(&track); err != nil {
		return models.Track{}, err
	}

	return track, nil
}

func (cache *CacheManager) StoreTrack(track models.Track) error {
	writer, err := os.Create(cache.InfoPath(track.ID))
	defer writer.Close()
	if err != nil {
		return err
	}

	storedTrack := track
	storedTrack.Message = ""
	storedTrack.User = ""
	storedTrack.Lang = ""

	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(storedTrack); err != nil {
		return err
	}

	return nil
}

func (cache *CacheManager) IsCached(trackID string) bool {
	return isFileExists(cache.InfoPath(trackID)) && isFileExists(cache.TrackPath(trackID))
}

func (cache *CacheManager) InfoPath(trackId string) string {
	return cache.getPath(trackId + cache.ExtInfo)
}

func (cache *CacheManager) TrackPath(trackId string) string {
	return cache.getPath(trackId + cache.ExtTrack)
}

func (cache *CacheManager) RemoveOldies() {
	size, _ := dirSize(cache.CachePath)

	cleanUpSize := int64(math.Max(0, float64(size/1000/1000)-float64(cache.MaxSize)))

	if cleanUpSize > 0 {
		files, _ := ioutil.ReadDir(cache.CachePath)
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Unix() < files[j].ModTime().Unix()
		})

		for _, file := range files {
			if strings.HasSuffix(file.Name(), cache.ExtTrack) {
				trackID := file.Name()[:len(file.Name())-len(cache.ExtTrack)]
				cleanUpSize -= file.Size() // the info file is small anyway :)
				os.Remove(cache.TrackPath(trackID))
				os.Remove(cache.InfoPath(trackID))

				if cleanUpSize <= 0 {
					break
				}
			}
		}
	}
}

func (cache *CacheManager) getPath(filename string) string {
	return path.Join(cache.CachePath, filename)
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func isFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	now := time.Now()
	os.Chtimes(path, now, now)
	return true
}

func cachePath() string {
	if cachePath, err := os.UserCacheDir(); err == nil {
		return path.Join(cachePath, "zukebox", "cache")
	}
	return path.Join("tmp", "zukebox", "cache")
}

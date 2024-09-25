package cache

import (
	"encoding/json"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"os"
	"path/filepath"
)

type Cache struct {
	directory   string
	trackUsages map[string]int
}

type Reader interface {
	TrackUsages(trackId string) int
}

func New(configDirectory string) *Cache {
	return &Cache{directory: filepath.Join(configDirectory, "cache")}
}

func (c *Cache) Initialized() bool {
	_, err := os.Stat(c.directory)
	return err == nil
}

func (c *Cache) Build(tracks []shared.Track, entries []shared.JournalEntry) error {
	err := os.RemoveAll(c.directory)
	if err != nil {
		return err
	}
	err = os.MkdirAll(c.directory, 0755)
	if err != nil {
		return err
	}

	usages := make(map[string]int)
	for _, entry := range entries {
		usages[entry.TrackId] = usages[entry.TrackId] + 1
	}
	payload, _ := json.Marshal(usages)
	return os.WriteFile(filepath.Join(c.directory, "trackUsages.json"), payload, 0644)
}

func (c *Cache) TrackUsages(trackId string) int {
	return c.trackUsages[trackId]
}

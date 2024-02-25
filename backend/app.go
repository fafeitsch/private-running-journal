package backend

import (
	"context"
	"fmt"
	"github.com/twpayne/go-gpx"
	"log"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time now!", name)
}

type JournalListEntry struct {
	Id           string `json:"id"`
	Date         string `json:"date"`
	TrackName    string `json:"trackName"`
	TrackVariant string `json:"trackVariant"`
	Length       int    `json:"length"`
}

type JournalEntry struct {
	Id    string `json:"id"`
	Date  string `json:"date"`
	Track Track  `json:"track"`
}

type Track struct {
	Id           string `json:"id"`
	Length       int    `json:"length"`
	Name         string `json:"name"`
	Variant      string `json:"variant"`
	WaypointData string `json:"waypointData"`
	Comment      string `json:"comment"`
}

type TrackListEntry struct {
	Name   string  `json:"name"`
	Length float64 `json:"length"`
}

func (a *App) GetJournalListEntries() []JournalListEntry {
	return make([]JournalListEntry, 0)
}

func (a *App) GetJournalEntry() JournalEntry {
	return JournalEntry{}
}

func (a *App) GetTracks() []TrackListEntry {
	log.Printf("getting tracks")
	files, err := os.ReadDir("appdata/geojson")
	if err != nil {
		log.Printf("could not read appdata/geojson: %v", err)
	}

	result := make([]TrackListEntry, 0, len(files))
	for _, file := range files {
		path := fmt.Sprintf("./appdata/geojson/%s", file.Name())
		fileContent, err := os.Open(path)
		if err != nil {
			log.Printf("could not read file %s: %v", path, err)
		}
		tracks, err := gpx.Read(fileContent)
		if err != nil {
			log.Printf("could not parse file %s: %v", file, err)
		}
		for _, track := range tracks.Trk {
			result = append(
				result, TrackListEntry{Name: track.Name, Length: distance(track.TrkSeg[0].TrkPt)},
			)
		}
	}
	return result
}

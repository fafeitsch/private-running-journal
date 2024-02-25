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
	Id            string `json:"id"`
	Date          string `json:"date"`
	TrackBaseName string `json:"trackBaseName"`
	TrackVariant  string `json:"trackVariant"`
	Length        int    `json:"length"`
}

type JournalEntry struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Track   Track  `json:"track"`
	Comment string `json:"comment"`
	Time    string `json:"time"`
}

type Track struct {
	Id           string `json:"id"`
	Length       int    `json:"length"`
	BaseName     string `json:"baseName"`
	BaseId       string `json:"baseId"`
	Variant      string `json:"variant"`
	WaypointData string `json:"waypointData"`
}

type TrackListEntryOld struct {
	Name   string  `json:"name"`
	Length float64 `json:"length"`
}

func (a *App) GetJournalListEntries() []JournalListEntry {
	return make([]JournalListEntry, 0)
}

func (a *App) GetJournalEntry() JournalEntry {
	return JournalEntry{}
}

func (a *App) GetTracks() []TrackListEntryOld {
	log.Printf("getting tracks")
	files, err := os.ReadDir("appdata/geojson")
	if err != nil {
		log.Printf("could not read appdata/geojson: %v", err)
	}

	result := make([]TrackListEntryOld, 0, len(files))
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
				result, TrackListEntryOld{Name: track.Name, Length: distance(track.TrkSeg[0].TrkPt)},
			)
		}
	}
	return result
}

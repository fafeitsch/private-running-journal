package backend

import (
	"context"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/tracks"
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
	Id      string       `json:"id"`
	Date    string       `json:"date"`
	Track   tracks.Track `json:"track"`
	Comment string       `json:"comment"`
	Time    string       `json:"time"`
}

func (a *App) GetJournalListEntries() []JournalListEntry {
	return make([]JournalListEntry, 0)
}

func (a *App) GetJournalEntry() JournalEntry {
	return JournalEntry{}
}

func (a *App) GetTracks() ([]tracks.Track, error) {
	return tracks.GetTracks("appdata")
}

func (a *App) GetGpxData(baseName string, variant string) (tracks.GpxData, error) {
	return tracks.GetGpxData("appdata", baseName, variant)
}

package backend

import (
	"context"
	"github.com/fafeitsch/local-track-journal/backend/journal"
	"github.com/fafeitsch/local-track-journal/backend/tracks"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	err := tracks.Init("appdata")
	if err != nil {
		log.Fatalf("could not initialize track directory: %v", err)
	}
}

func (a *App) GetJournalListEntries() ([]journal.ListEntry, error) {
	return journal.ReadEntries("appdata")
}

func (a *App) GetJournalEntry(id string) (journal.Entry, error) {
	return journal.ReadJournalEntry("appdata", id)
}

func (a *App) CreateJournalEntry(date string, trackId string) (journal.ListEntry, error) {
	return journal.CreateEntry("appdata", date, trackId)
}

func (a *App) GetTracks() []tracks.Track {
	return tracks.RootTracks()
}

func (a *App) GetGpxData(id string) (tracks.GpxData, error) {
	return tracks.GetGpxData("appdata/tracks", id)
}

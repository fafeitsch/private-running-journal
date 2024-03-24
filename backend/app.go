package backend

import (
	"context"
	"github.com/fafeitsch/local-track-journal/backend/httpapi"
	"github.com/fafeitsch/local-track-journal/backend/journal"
	"github.com/fafeitsch/local-track-journal/backend/tracks"
	"log"
	"net/http"
	"sync"
)

// App struct
type App struct {
	ctx     context.Context
	tracks  *tracks.Tracks
	journal *journal.Journal
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	var err error
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		a.tracks, err = tracks.New("appdata")
		if err != nil {
			log.Fatalf("could not initialize track directory: %v", err)
		}
		group.Done()
	}()
	go func() {
		a.journal, err = journal.New("appdata")
		if err != nil {
			log.Fatalf("could not initialize journal: %v", err)
		}
		group.Done()
	}()
	tileServer := httpapi.NewTileServer("appdata", "https://tile.openstreetmap.org/{z}/{x}/{y}.png", true)
	go func() {
		err = http.ListenAndServe("127.0.0.1:47836", tileServer)
		if err != nil {
			log.Fatalf("could not start tile server: %v", err)
		}
	}()
	group.Wait()
}

func (a *App) GetJournalListEntries() ([]journal.ListEntry, error) {
	return a.journal.ReadEntries()
}

func (a *App) GetJournalEntry(id string) (journal.Entry, error) {
	return a.journal.ReadJournalEntry(id)
}

func (a *App) SaveJournalEntry(entry journal.Entry) error {
	return a.journal.SaveEntry(entry)
}

func (a *App) CreateJournalEntry(date string, trackId string) (journal.ListEntry, error) {
	return a.journal.CreateEntry(date, trackId)
}

func (a *App) GetTracks() []tracks.Track {
	return a.tracks.RootTracks()
}

func (a *App) GetGpxData(id string) (tracks.GpxData, error) {
	return tracks.GetGpxData("appdata/tracks", id)
}

func (a *App) ComputePolylineProps(coords []tracks.Coordinates) tracks.PolylineProps {
	return tracks.ComputePolylineProps(coords)
}

func (a *App) CreateNewTrack(track tracks.CreateTrack) (*tracks.Track, error) {
	return a.tracks.CreateTrack(track)
}

func (a *App) SaveTrack(track tracks.SaveTrack) (*tracks.Track, error) {
	return a.tracks.SaveTrack(track)
}

package backend

import (
	"context"
	"github.com/fafeitsch/local-track-journal/backend/httpapi"
	"github.com/fafeitsch/local-track-journal/backend/journal"
	"github.com/fafeitsch/local-track-journal/backend/settings"
	"github.com/fafeitsch/local-track-journal/backend/tracks"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// App struct
type App struct {
	ctx             context.Context
	configDirectory string
	tracks          *tracks.Tracks
	journal         *journal.Journal
	settings        *settings.Settings
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.setupConfigDirectory()
	a.ctx = ctx
	var err error
	a.settings, err = settings.New(a.configDirectory)
	if err != nil {
		log.Fatalf("could not read settings: %v", err)
	}
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		a.tracks, err = tracks.New(a.configDirectory)
		if err != nil {
			log.Fatalf("could not initialize track directory: %v", err)
		}
		group.Done()
	}()
	go func() {
		a.journal, err = journal.New(a.configDirectory)
		if err != nil {
			log.Fatalf("could not initialize journal: %v", err)
		}
		group.Done()
	}()
	group.Wait()
	tileServer := httpapi.NewTileServer(
		a.configDirectory, a.settings.MapSettings().TileServer, a.settings.MapSettings().CacheTiles,
	)
	go func() {
		err = http.ListenAndServe("127.0.0.1:47836", tileServer)
		if err != nil {
			log.Fatalf("could not start tile server: %v", err)
		}
	}()
}

func (a *App) setupConfigDirectory() {
	var homeDir, homeDirErr = os.UserHomeDir()
	if len(os.Args) > 1 {
		a.configDirectory = os.Args[1]
	} else if homeDirErr == nil {
		a.configDirectory = filepath.Join(homeDir, ".local-track-journal")
	} else {
		log.Fatalf(
			"cannot read home directory, please specify one by providing a command line argument: %v", homeDirErr,
		)
	}
	err := os.MkdirAll(a.configDirectory, os.ModePerm)
	if err != nil {
		log.Fatalf("could not create app's config dir: %v", err)
	}
	log.Printf("setting app's home dir to %s", a.configDirectory)
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

func (a *App) DeleteJournalEntry(id string) error {
	return a.journal.DeleteEntry(id)
}

func (a *App) GetTracks() []tracks.Track {
	return a.tracks.RootTracks()
}

func (a *App) GetGpxData(id string) (tracks.GpxData, error) {
	return a.tracks.GetGpxData(id)
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

func (a *App) GetSettings() settings.AppSettings {
	return a.settings.AppSettings()
}

func (a *App) SaveSettings(settings settings.AppSettings) error {
	return a.settings.SaveSettings(settings)
}

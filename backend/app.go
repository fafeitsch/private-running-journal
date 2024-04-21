package backend

import (
	"context"
	"github.com/fafeitsch/private-running-journal/backend/backup"
	"github.com/fafeitsch/private-running-journal/backend/httpapi"
	"github.com/fafeitsch/private-running-journal/backend/journal"
	"github.com/fafeitsch/private-running-journal/backend/settings"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/fafeitsch/private-running-journal/backend/tracks"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type App struct {
	ctx             context.Context
	configDirectory string
	tracks          *tracks.Tracks
	journal         *journal.Journal
	settings        *settings.Settings
	backup          *backup.Backup
}

func NewApp() *App {
	a := &App{}
	a.setupConfigDirectory()
	var err error
	a.settings, err = settings.New(a.configDirectory)
	if err != nil {
		log.Fatalf("could not read settings: %v", err)
	}
	a.backup = backup.Init(
		a.configDirectory, a.settings.GitSettings().Enabled, a.settings.GitSettings().PushAfterCommit,
	)
	if a.settings.GitSettings().Enabled && a.settings.GitSettings().PullOnStartUp {
		log.Printf("pulling")
		err = a.backup.Pull()
		if err != nil {
			log.Fatalf("could not pull: %v", err)
		}
	}
	return a
}

func (a *App) HeadlessMode() bool {
	return a.settings.AppSettings().HeadlessMode
}

func (a *App) Language() string {
	return a.settings.AppSettings().Language
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	shared.Context = ctx
	var err error
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
	log.Printf("use headless mode: %v", a.HeadlessMode())
	log.Printf("start up")
}

func (a *App) setupConfigDirectory() {
	var homeDir, homeDirErr = os.UserHomeDir()
	if len(os.Args) > 1 {
		a.configDirectory = os.Args[1]
	} else if homeDirErr == nil {
		a.configDirectory = filepath.Join(homeDir, ".private-running-journal")
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
	return a.tracks.Tracks()
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

func (a *App) DeleteTrack(id string) error {
	return a.tracks.DeleteTrack(id)
}

func (a *App) MoveTrack(id string, newPath string) (*tracks.Track, error) {
	return a.tracks.MoveTrack(id, newPath)
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

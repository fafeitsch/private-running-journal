package backend

import (
	"context"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/backup"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/httpapi"
	"github.com/fafeitsch/private-running-journal/backend/journal"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/settings"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/fafeitsch/private-running-journal/backend/tracks"
	"github.com/fafeitsch/private-running-journal/backend/tracks/trackEditor"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	ctx             context.Context
	configDirectory string
	tracks          *tracks.Tracks
	trackEditor     *trackEditor.TrackEditor
	journal         *journal.Journal
	settings        *settings.Settings
	backup          *backup.Backup
	cache           *projection.Projection
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

	service := filebased.NewService(a.configDirectory)
	a.tracks = tracks.New(a.configDirectory, service)
	a.journal, err = journal.New(a.configDirectory, service)
	if err != nil {
		log.Fatalf("could not initialize journal: %v", err)
	}
	trackUsagesProjector := projection.TrackUsagesProjector{Journal: a.journal}
	a.trackEditor = trackEditor.New(service, &trackUsagesProjector)
	projectors := make([]projection.Projector, 0)
	projectors = append(projectors, &trackUsagesProjector)
	for i := range a.tracks.Projectors {
		projectors = append(projectors, a.tracks.Projectors[i])
	}
	a.cache = projection.New(a.configDirectory, projectors...)
	if !a.cache.Initialized() {
		err = a.cache.Build()
		if err != nil {
			log.Fatalf("could not initialize projections: %v", err)
		}
	}
	a.tracks.TrackUsagesProjector = &trackUsagesProjector

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

func (a *App) GetJournalListEntries(start *string, end *string) ([]journal.ListEntry, error) {
	var startDate *time.Time
	if start != nil {
		date, err := time.Parse(time.RFC3339, *start)
		if err != nil {
			return []journal.ListEntry{}, fmt.Errorf("could not parse start date: %v", err)
		}
		startDate = &date
	}
	var endDate *time.Time
	if end != nil {
		date, err := time.Parse(time.RFC3339, *end)
		if err != nil {
			return []journal.ListEntry{}, fmt.Errorf("could not parse end date: %v", err)
		}
		endDate = &date
	}
	return a.journal.ListEntries(startDate, endDate), nil
}

func (a *App) GetJournalEntry(id string) (journal.Entry, error) {
	return a.journal.ReadJournalEntry(id)
}

func (a *App) SaveJournalEntry(entry journal.Entry) error {
	return a.journal.SaveEntry(entry)
}

func (a *App) CreateJournalEntry(entry journal.Entry) (journal.ListEntry, error) {
	return a.journal.CreateEntry(entry)
}

func (a *App) DeleteJournalEntry(id string) error {
	return a.journal.DeleteEntry(id)
}

func (a *App) GetTracks() ([]tracks.TrackListEntry, error) {
	return a.tracks.Tracks()
}

func (a *App) GetTrackTree() (tracks.TrackTreeNode, error) {
	return a.tracks.TrackTree()
}

func (a *App) TrackEditor() *trackEditor.TrackEditor {
	return a.trackEditor
}

func (a *App) GetTrack(id string) (tracks.Track, error) {
	return a.tracks.GetTrack(id)
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

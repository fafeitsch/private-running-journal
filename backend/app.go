package backend

import (
	"context"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/application/journalEditor"
	"github.com/fafeitsch/private-running-journal/backend/application/journalList"
	"github.com/fafeitsch/private-running-journal/backend/application/trackEditor"
	"github.com/fafeitsch/private-running-journal/backend/backup"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/httpapi"
	"github.com/fafeitsch/private-running-journal/backend/journal"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/settings"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	ctx             context.Context
	configDirectory string
	trackEditor     *trackEditor.TrackEditor
	journalEditor   *journalEditor.JournalEditor
	journalList     *journalList.JournalList
	journal         *journal.Journal
	settings        *settings.Settings
	backup          *backup.Backup
	cache           *projection.Projection
	trackTree       *projection.TrackTree
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
	a.journal, err = journal.New(a.configDirectory, service)
	if err != nil {
		log.Fatalf("could not initialize journal: %v", err)
	}
	trackUsagesProjector := &projection.TrackUsages{}
	sortedJournalProjector := &projection.SortedJournalEntries{Directory: a.configDirectory}
	a.trackTree = &projection.TrackTree{}
	a.journalEditor = journalEditor.New(service)
	a.trackEditor = trackEditor.New(service, trackUsagesProjector)
	a.journalList = journalList.New(service, sortedJournalProjector)
	projectors := make([]projection.Projector, 0)
	projectors = append(projectors, trackUsagesProjector)
	projectors = append(projectors, a.trackTree)
	projectors = append(projectors, sortedJournalProjector)
	a.cache = projection.New(filepath.Join(a.configDirectory, ".projection"), service, projectors...)
	err = a.cache.Build()
	if err != nil {
		log.Fatalf("could not initialize projections: %v", err)
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
	log.Printf("program args: %v", os.Args)
	if os.Args[0] == "/tmp/wailsbindings" {
		a.configDirectory = "/tmp"
	} else if len(os.Args) > 1 {
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

func (a *App) GetJournalListEntries(start string, end string) ([]journalList.ListEntryDto, error) {
	var startDate time.Time
	date, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return []journalList.ListEntryDto{}, fmt.Errorf("could not parse start date: %v", err)
	}
	startDate = date
	var endDate time.Time
	date, err = time.Parse(time.RFC3339, end)
	if err != nil {
		return []journalList.ListEntryDto{}, fmt.Errorf("could not parse end date: %v", err)
	}
	endDate = date
	return a.journalList.ReadListEntries(startDate, endDate)
}

func (a *App) GetJournalEntry(id string) (journal.Entry, error) {
	return a.journal.ReadJournalEntry(id)
}

func (a *App) DeleteJournalEntry(id string) error {
	return a.journal.DeleteEntry(id)
}

func (a *App) GetTrackTree() projection.TrackTreeNode {
	return a.trackTree.Get()
}

func (a *App) TrackEditor() *trackEditor.TrackEditor {
	return a.trackEditor
}

func (a *App) JournalEditor() *journalEditor.JournalEditor {
	return a.journalEditor
}

func (a *App) GetSettings() settings.AppSettings {
	return a.settings.AppSettings()
}

func (a *App) SaveSettings(settings settings.AppSettings) error {
	return a.settings.SaveSettings(settings)
}

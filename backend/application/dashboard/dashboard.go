package dashboard

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"time"
)

type DashboardDto struct {
	TotalDistance int     `json:"totalDistance"`
	TopTracks     []Track `json:"topTracks"`
}

type Track struct {
	Name    string
	Id      string
	Parents []string
}

type Assembler struct {
	sortedEntries *projection.SortedJournalEntries
	fileService   *filebased.Service
}

func NewAssembler(sortedEntries *projection.SortedJournalEntries, fileService *filebased.Service) *Assembler {
	return &Assembler{sortedEntries: sortedEntries, fileService: fileService}
}

type Options struct {
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	TopTracks int       `json:"topTracks"`
}

func (a *Assembler) LoadDashboard(options Options) (*DashboardDto, error) {
	entries, err := a.sortedEntries.FindJournalEntryIdsBetween(options.From, options.To)
	if err != nil {
		return nil, err
	}
	sum := 0
	trackCache := make(map[string]shared.Track)
	for _, entry := range entries {
		loaded, err := a.fileService.ReadJournalEntry(entry)
		if err != nil {
			return nil, err
		}
		if loaded.CustomLength != nil {
			sum = sum + *loaded.CustomLength
			continue
		}
		track, ok := trackCache[loaded.TrackId]
		if !ok {
			track, err = a.fileService.ReadTrack(loaded.TrackId)
			if err != nil {
				return nil, err
			}
			trackCache[loaded.TrackId] = track
		}
		sum = sum + (track.Waypoints.Length() * loaded.Laps)
	}
	return &DashboardDto{
		TotalDistance: sum,
		TopTracks:     []Track{},
	}, nil
}

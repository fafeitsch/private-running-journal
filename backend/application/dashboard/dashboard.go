package dashboard

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"math"
	"slices"
	"time"
)

type DashboardDto struct {
	TotalDistance int     `json:"totalDistance"`
	TopTracks     []Track `json:"topTracks"`
	TotalRuns     int     `json:"totalRuns"`
}

type Track struct {
	Name    string   `json:"name"`
	Id      string   `json:"id"`
	Parents []string `json:"parents"`
	Count   int      `json:"count"`
	Length  int      `json:"length"`
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
	trackCounter := make(map[string]int)
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
		trackCounter[loaded.TrackId] = trackCounter[loaded.TrackId] + 1
		if !ok {
			track, err = a.fileService.ReadTrack(loaded.TrackId)
			if err != nil {
				return nil, err
			}
			trackCache[loaded.TrackId] = track
		}
		sum = sum + (track.Waypoints.Length() * loaded.Laps)
	}
	topTracks := make([]Track, 0, options.TopTracks)
	for id, track := range trackCache {
		topTracks = append(topTracks, Track{Id: id, Name: track.Name, Count: trackCounter[id], Parents: track.Parents, Length: track.Waypoints.Length()})
	}
	slices.SortFunc(topTracks, func(e Track, e2 Track) int {
		return e2.Count - e.Count
	})
	return &DashboardDto{
		TotalDistance: sum,
		TopTracks:     topTracks[:int(math.Min(float64(options.TopTracks), float64(len(topTracks))))],
		TotalRuns:     len(entries),
	}, nil
}

package journalList

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"path/filepath"
	"time"
)

type JournalList struct {
	fileService *filebased.Service
	trackLookup *projection.TrackLookup
}

func New(service *filebased.Service, trackLookup *projection.TrackLookup) *JournalList {
	return &JournalList{
		fileService: service,
		trackLookup: trackLookup,
	}
}

type ListEntryDto struct {
	TrackName  string `json:"trackName"`
	TrackError bool   `json:"trackError"`
	Length     int    `json:"length"`
	Date       string `json:"date"`
	Id         string `json:"id"`
}

func (j *JournalList) ReadListEntries(start time.Time, end time.Time) ([]ListEntryDto, error) {
	result := make([]ListEntryDto, 0)
	files, err := j.fileService.ReadJournalEntriesBetween(start, end)
	if err != nil {
		return nil, fmt.Errorf("error reading journal entries: %v", err)
	}
	trackCache := make(map[string]shared.Track)
	for _, file := range files {
		track, ok := trackCache[file.TrackId]
		entry := ListEntryDto{Id: file.Id, Date: file.Date.Format(time.DateOnly)}
		if !ok {
			path := j.trackLookup.Get()[file.TrackId]
			track, err = j.fileService.ReadTrack(filepath.Join(path...))
			if err != nil {
				entry.TrackName = file.TrackId
				entry.TrackError = true
				log.Printf("could not read track of joural entry %s: %v", file.Id, err)
			}
		}
		entry.Length = track.Waypoints.Length() * file.Laps
		if file.CustomLength > 0 {
			entry.Length = file.CustomLength
		}
		entry.TrackName = track.Name
		trackCache[file.TrackId] = track
		result = append(result, entry)
	}
	return result, nil
}

package journalList

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"time"
)

type JournalList struct {
	fileService            *filebased.Service
	sortedJournalProjector *projection.SortedJournalEntries
}

func New(service *filebased.Service, sortedJournalProjector *projection.SortedJournalEntries) *JournalList {
	return &JournalList{
		fileService:            service,
		sortedJournalProjector: sortedJournalProjector,
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
	ids, err := j.sortedJournalProjector.FindJournalEntryIdsBetween(start, end)
	if err != nil {
		return nil, fmt.Errorf("error reading journal entries: %v", err)
	}
	trackCache := make(map[string]shared.Track)
	for _, journalId := range ids {
		file, err := j.fileService.ReadJournalEntry(journalId)
		if err != nil {
			log.Printf("could not read journal entry with id \"%s\": %v", journalId, err)
			continue
		}
		track, ok := trackCache[file.TrackId]
		entry := ListEntryDto{Id: file.Id, Date: file.Date.Format(time.DateOnly)}
		if !ok {
			track, err = j.fileService.ReadTrack(file.TrackId)
			if err != nil {
				entry.TrackName = file.TrackId
				entry.TrackError = true
				log.Printf("could not read track of joural entry %s: %v", file.Id, err)
			}
		}
		entry.Length = track.Waypoints.Length() * file.Laps
		if file.CustomLength != nil {
			entry.Length = *file.CustomLength
		}
		entry.TrackName = track.Name
		trackCache[file.TrackId] = track
		result = append(result, entry)
	}
	return result, nil
}

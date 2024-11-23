package journalEditor

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"time"
)

type JournalEditor struct {
	fileService *filebased.Service
	trackLookup *projection.TrackLookup
}

type SaveEntryDto struct {
	Id           string `json:"id"`
	TrackId      string `json:"trackId"`
	Date         string `json:"date"`
	Comment      string `json:"comment"`
	Time         string `json:"time"`
	Laps         int    `json:"laps"`
	CustomLength int    `json:"customLength"`
}

func New(service *filebased.Service, trackLookup *projection.TrackLookup) *JournalEditor {
	return &JournalEditor{fileService: service, trackLookup: trackLookup}
}

type SaveJournalEntryResultDto struct {
	Id string `json:"id"`
}

func (j *JournalEditor) SaveJournalEntry(entry SaveEntryDto) (SaveJournalEntryResultDto, error) {
	oldTrackId := ""
	existing, err := j.fileService.ReadJournalEntry(entry.Id)
	if err == nil {
		oldTrackId = existing.TrackId
	}
	date, _ := time.Parse(time.DateOnly, entry.Date)
	journalEntry := shared.JournalEntry{
		TrackId:      entry.TrackId,
		Id:           entry.Id,
		Date:         date,
		Comment:      entry.Comment,
		CustomLength: entry.CustomLength,
		Laps:         entry.Laps,
		Time:         entry.Time,
	}
	id, err := j.fileService.SaveJournalEntry(
		journalEntry,
	)
	if err != nil {
		return SaveJournalEntryResultDto{}, fmt.Errorf("could not write journal entry: %v", err)
	}
	shared.SendEvent(
		shared.JournalEntryUpsertedEvent{
			JournalEntry: &journalEntry,
			OldTrackId:   oldTrackId,
		},
	)
	return SaveJournalEntryResultDto{Id: id}, nil
}

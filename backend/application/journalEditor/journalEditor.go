package journalEditor

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"time"
)

type JournalEditor struct {
	fileService *filebased.Service
}

type SaveEntryDto struct {
	Id           string `json:"id"`
	TrackId      string `json:"trackId"`
	Date         string `json:"date"`
	Comment      string `json:"comment"`
	Time         string `json:"time"`
	Laps         int    `json:"laps"`
	CustomLength *int   `json:"customLength"`
}

type EntryDto struct {
	Id           string `json:"id"`
	TrackId      string `json:"trackId"`
	Date         string `json:"date"`
	Comment      string `json:"comment"`
	Time         string `json:"time"`
	Laps         int    `json:"laps"`
	CustomLength *int   `json:"customLength"`
}

func New(service *filebased.Service) *JournalEditor {
	return &JournalEditor{fileService: service}
}

type SaveJournalEntryResultDto struct {
	Id string `json:"id"`
}

func (j *JournalEditor) GetJournalEntry(id string) (EntryDto, error) {
	existing, err := j.fileService.ReadJournalEntry(id)
	if err != nil {
		return EntryDto{}, fmt.Errorf("could not read journal entry: %v", err)
	}
	return EntryDto{
		Id:           existing.Id,
		TrackId:      existing.TrackId,
		Date:         existing.Date.Format(time.DateOnly),
		Comment:      existing.Comment,
		Time:         existing.Time,
		Laps:         existing.Laps,
		CustomLength: existing.CustomLength,
	}, nil
}

func (j *JournalEditor) SaveJournalEntry(entry SaveEntryDto) (SaveJournalEntryResultDto, error) {
	oldTrackId := ""
	var oldDate *time.Time
	existing, err := j.fileService.ReadJournalEntry(entry.Id)
	if err == nil {
		oldTrackId = existing.TrackId
		oldDate = &existing.Date
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
	err = j.fileService.SaveJournalEntry(
		journalEntry,
	)
	if err != nil {
		return SaveJournalEntryResultDto{}, fmt.Errorf("could not write journal entry: %v", err)
	}
	shared.SendEvent(
		shared.JournalEntryUpsertedEvent{
			JournalEntry: &journalEntry,
			OldTrackId:   oldTrackId,
			OldDate:      oldDate,
		},
	)
	return SaveJournalEntryResultDto{Id: entry.Id}, nil
}

func (j *JournalEditor) DeleteJournalEntry(id string) error {
	existing, err := j.fileService.ReadJournalEntry(id)
	if err == nil {
		shared.SendEvent(shared.JournalEntryDeletedEvent{JournalEntry: &existing})
	}
	return j.fileService.DeleteJournalEntry(id)
}

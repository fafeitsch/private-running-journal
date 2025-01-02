package filebased

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"time"
)

type entryFile struct {
	Id           string `json:"id"`
	Track        string `json:"track"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Comment      string `json:"comment"`
	Laps         int    `json:"laps"`
	CustomLength *int   `json:"customLength,omitempty"`
}

func (s *Service) ReadAllJournalEntries() ([]shared.JournalEntry, error) {
	result := make([]shared.JournalEntry, 0, 0)
	walkPath := filepath.Join(s.path, journalDirectory)
	err := filepath.Walk(
		walkPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			if path == walkPath || info.IsDir() {
				return nil
			}

			if info.Name() != "entry.json" {
				return filepath.SkipDir
			}
			listEntry, err := s.ReadJournalEntry(filepath.Base(filepath.Dir(path)))
			if err != nil {
				return err
			}
			result = append(result, listEntry)
			return nil
		},
	)
	return result, err
}

func (s *Service) ReadJournalEntry(id string) (shared.JournalEntry, error) {
	var listEntry entryFile
	file, err := os.Open(filepath.Join(s.path, journalDirectory, id[0:2], id, "entry.json"))
	if err != nil {
		return shared.JournalEntry{}, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&listEntry)
	if err != nil {
		return shared.JournalEntry{}, fmt.Errorf("could not parse file: %v", err)
	}
	date, err := time.Parse(time.DateOnly, listEntry.Date)
	if err != nil {
		return shared.JournalEntry{}, fmt.Errorf("could not parse date: %v", err)
	}
	var customLength *int
	if listEntry.CustomLength != nil {
		customLength = listEntry.CustomLength
	}
	return shared.JournalEntry{
		TrackId:      listEntry.Track,
		Id:           id,
		Date:         date,
		Comment:      listEntry.Comment,
		CustomLength: customLength,
		Laps:         listEntry.Laps,
		Time:         listEntry.Time,
	}, nil
}

func (s *Service) SaveJournalEntry(entry shared.JournalEntry) error {
	log.Printf("save journal entry: %v", entry)
	path := filepath.Join(s.path, journalDirectory, entry.Id[0:2], entry.Id)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("could not create directory: %v", err)
	}
	payload, _ := json.Marshal(
		entryFile{
			Id:           entry.Id,
			Track:        entry.TrackId,
			Laps:         entry.Laps,
			Date:         entry.Date.Format(time.DateOnly),
			Time:         entry.Time,
			Comment:      entry.Comment,
			CustomLength: entry.CustomLength,
		},
	)
	return os.WriteFile(filepath.Join(path, "entry.json"), payload, 0644)
}

func (s *Service) DeleteJournalEntry(id string) error {
	return os.RemoveAll(filepath.Join(s.path, journalDirectory, id[0:2], id))
}

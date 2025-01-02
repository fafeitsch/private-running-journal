package projection

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const sortedJournalEntriesDirectory = "sortedJournalEntries"

type SortedJournalEntries struct {
	fileService *filebased.Service
	Directory   string
}

func (s *SortedJournalEntries) ProjectionName() string {
	return "sortedJournalEntries"
}

func (s *SortedJournalEntries) Init(message json.RawMessage, writer func()) {
	shared.Listen(
		shared.JournalEntryUpsertedEvent{}, func(k shared.JournalEntryUpsertedEvent) {
			s.handleUpsertEvent(k)
		},
	)
	shared.Listen(shared.JournalEntryDeletedEvent{}, func(k shared.JournalEntryDeletedEvent) {
		s.handleDeleteEvent(k.Date, k.Id)
	})
}

func (s *SortedJournalEntries) AddTrack(track shared.Track) {
}

func (s *SortedJournalEntries) AddJournalEntry(entry shared.JournalEntry) {
	s.handleUpsertEvent(shared.JournalEntryUpsertedEvent{
		JournalEntry: &entry,
		OldTrackId:   "",
		OldDate:      nil,
	})
}

func (s *SortedJournalEntries) Get() map[string]string {
	return make(map[string]string)
}

func (s *SortedJournalEntries) GetData() any {
	return make(map[string]string)
}

func (s *SortedJournalEntries) handleDeleteEvent(date time.Time, id string) {
	split := strings.Split(date.Format(time.DateOnly), "-")
	path := filepath.Join(s.Directory, sortedJournalEntriesDirectory, split[0], split[1], split[2], id)
	err := os.Remove(path)
	if err != nil {
		log.Printf("error removing file %s: %s", path, err)
	}
}

func (s *SortedJournalEntries) handleUpsertEvent(event shared.JournalEntryUpsertedEvent) {
	fmt.Printf("handle upsert event")
	if event.OldDate != nil {
		s.handleDeleteEvent(event.Date, event.Id)
	}
	split := strings.Split(event.Date.Format(time.DateOnly), "-")
	path := filepath.Join(s.Directory, ".projection", sortedJournalEntriesDirectory, split[0], split[1], split[2])
	os.MkdirAll(path, 0755)
	err := os.Symlink(filepath.Join("..", "..", "..", "..", "..", "journal", event.Id[0:2], event.Id), filepath.Join(path, event.Id))
	if err != nil {
		log.Printf("could not link journal entry: %v", err)
	}
}

var yearMatcher = regexp.MustCompile(`^\d{4}$`)
var monthMatcher = regexp.MustCompile(`^\d{2}$`)
var dayMatcher = regexp.MustCompile(`^\d{2}$`)

func (s *SortedJournalEntries) FindJournalEntryIdsBetween(start time.Time, end time.Time) ([]string, error) {
	result := make([]string, 0, 0)
	walkPath := filepath.Join(s.Directory, ".projection", sortedJournalEntriesDirectory)
	err := filepath.Walk(
		walkPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			parts := strings.Split(path, string(os.PathSeparator))
			partsLength := len(parts)
			if path == walkPath {
				return nil
			}
			if yearMatcher.MatchString(parts[partsLength-1]) {
				// we are on level one: years
				year, _ := strconv.Atoi(parts[partsLength-1])
				if year >= start.Year() && year <= end.Year() {
					return nil
				}
			} else if yearMatcher.MatchString(parts[partsLength-2]) && monthMatcher.MatchString(parts[partsLength-1]) {
				// we are on level three: months (no performance improvement here)
				return nil
			} else if yearMatcher.MatchString(parts[partsLength-3]) && monthMatcher.MatchString(parts[partsLength-2]) && dayMatcher.MatchString(parts[partsLength-1]) {
				// we are on level three: days
				year := parts[partsLength-3]
				month := parts[partsLength-2]
				day := parts[partsLength-1][:2]
				date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%s-%s-%s", year, month, day))
				if date.Equal(start) || (date.After(start) && date.Before(end)) {
					return nil
				} else {
					return filepath.SkipDir
				}
			} else if yearMatcher.MatchString(parts[partsLength-4]) && monthMatcher.MatchString(parts[partsLength-3]) && dayMatcher.MatchString(parts[partsLength-2]) {
				// below the days-level: this must be the id
				result = append(result, filepath.Base(path))
			}
			return nil
		},
	)
	return result, err
}

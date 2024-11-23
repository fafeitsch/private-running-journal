package filebased

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var yearMatcher = regexp.MustCompile(`\d{4}`)
var monthMatcher = regexp.MustCompile(`\d{2}`)
var dayMatcher = regexp.MustCompile(`\d{2}[a-zA-Z]*`)

type entryFile struct {
	Track        string `json:"track"`
	Time         string `json:"time"`
	Comment      string `json:"comment"`
	Laps         int    `json:"laps"`
	CustomLength int    `json:"customLength,omitempty"`
}

func (s *Service) ReadJournalEntriesBetween(start time.Time, end time.Time) ([]shared.JournalEntry, error) {
	result := make([]shared.JournalEntry, 0, 0)
	walkPath := filepath.Join(s.path, journalDirectory)
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
			if info.IsDir() {
				if yearMatcher.MatchString(parts[partsLength-1]) {
					// we are on level one: years
					year, _ := strconv.Atoi(parts[partsLength-1])
					if year >= start.Year() && year <= end.Year() {
						return nil
					}
				} else if yearMatcher.MatchString(parts[partsLength-2]) && monthMatcher.MatchString(parts[partsLength-2]) {
					// we are on leven three: months (no performance improvement here)
					return nil
				} else if yearMatcher.MatchString(parts[partsLength-3]) && monthMatcher.MatchString(parts[partsLength-2]) && monthMatcher.MatchString(parts[partsLength-1]) {
					// we are on level three: days
					year := parts[partsLength-3]
					month := parts[partsLength-2]
					day := parts[partsLength-1][:2]
					date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%s-%s-%s", year, month, day))
					if date.Equal(start) || (date.After(start) && date.Before(end)) {
						return nil
					}
				}
				return filepath.SkipDir
			}
			if info.Name() != "entry.json" {
				return filepath.SkipDir
			}
			var listEntry entryFile
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("could not open file: %v", err)
			}
			defer file.Close()
			err = json.NewDecoder(file).Decode(&listEntry)
			if err != nil {
				return fmt.Errorf("could not parse file: %v", err)
			}
			year := parts[partsLength-4]
			month := parts[partsLength-3]
			day := parts[partsLength-2][:2]
			date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%s-%s-%s", year, month, day))
			result = append(
				result, shared.JournalEntry{
					TrackId:      listEntry.Track,
					Id:           filepath.Join(parts[partsLength-4], parts[partsLength-3], parts[partsLength-2]),
					Date:         date,
					Comment:      listEntry.Comment,
					CustomLength: listEntry.CustomLength,
					Laps:         listEntry.Laps,
				},
			)
			return nil
		},
	)
	return result, err
}

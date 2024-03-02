package journal

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/tracks"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ListEntry struct {
	Id          string   `json:"id"`
	Date        string   `json:"date"`
	ParentNames []string `json:"trackParents"`
	TrackName   string   `json:"trackName"`
	Length      int      `json:"length"`
}

type entryFile struct {
	Track   string `json:"track"`
	Time    string `json:"time"`
	Comment string `json:"comment"`
}

var directoryRegex = regexp.MustCompile("\\d\\d[a-z]?|\\d\\d\\d\\d")

func ReadEntries(baseDirectory string) ([]ListEntry, error) {
	result := make([]ListEntry, 0, 0)
	journalDirectory := filepath.Join(baseDirectory, "journal")
	err := filepath.Walk(
		journalDirectory, func(path string, info os.FileInfo, err error) error {
			fmt.Printf("walk %s, journalDir: %s", path, journalDirectory)
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			if info.IsDir() || info.Name() != "entry.json" {
				if path == journalDirectory || directoryRegex.MatchString(info.Name()) {
					return nil
				}
				log.Printf(
					"skipping directory \"%s\" because it does not have the required format yyyy/mm/dd", path,
				)
				return filepath.SkipDir
			}
			entry, err := ReadJournalEntry(
				baseDirectory,
				strings.Replace(strings.Replace(path, journalDirectory, "", 1), info.Name(), "", 1),
			)
			if err != nil {
				log.Printf("skipping journal entry \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			listEntry := ListEntry{Id: entry.Id, Date: entry.Date}
			if entry.Track != nil {
				listEntry.ParentNames = entry.Track.ParentNames
				listEntry.TrackName = entry.Track.Name
				listEntry.Length = entry.Track.Length
			}
			result = append(result, listEntry)
			return nil
		},
	)
	return result, err
}

type Entry struct {
	Id      string        `json:"id"`
	Date    string        `json:"date"`
	Track   *tracks.Track `json:"track"`
	Comment string        `json:"comment"`
	Time    string        `json:"time"`
}

func ReadJournalEntry(basePath string, path string) (Entry, error) {
	file, err := os.Open(filepath.Join(basePath, "journal", path, "entry.json"))
	if err != nil {
		log.Printf("skipping journal entry \"%s\" because an error occurred: %v", path, err)
		return Entry{}, err
	}
	var entryDescriptor entryFile
	err = json.NewDecoder(file).Decode(&entryDescriptor)
	if err != nil {
		log.Printf("skipping journal entry \"%s\" because an error occurred: %v", path, err)
		return Entry{}, err
	}
	journalEntry := Entry{Id: path, Comment: entryDescriptor.Comment, Time: entryDescriptor.Time}
	date, err := computeDateFromPath(path)
	if err != nil {
		log.Printf("skipping journal entry \"%s\" because an error occurred: %v", path, err)
		return Entry{}, err
	}
	journalEntry.Date = date
	track, ok := tracks.GetTrack(entryDescriptor.Track)
	if ok {
		journalEntry.Track = &track
	}
	return journalEntry, nil
}

var pathRegex = regexp.MustCompile("(\\d\\d\\d\\d)/(\\d\\d)/(\\d\\d)[a-z]?/")

func computeDateFromPath(path string) (string, error) {
	regexResult := pathRegex.FindStringSubmatch(path)
	if regexResult == nil {
		return "", fmt.Errorf("path \"%s\" does not have the correct format dddd/dd/dd/*.json", path)
	}
	return fmt.Sprintf("%s-%s-%s", regexResult[1], regexResult[2], regexResult[3]), nil
}

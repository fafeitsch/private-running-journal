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
	Laps    int    `json:"laps"`
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
			listEntry, err := readListEntry(
				baseDirectory, strings.Replace(strings.Replace(path, journalDirectory, "", 1), info.Name(), "", 1),
			)
			if err != nil {
				log.Printf("skipping journal entry \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
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
	Laps    int           `json:"laps"`
	//needed in FE if the Track file is missing to give a reference what track should be there
	LinkedTrack string `json:"linkedTrack"`
}

func ReadJournalEntry(basePath string, path string) (Entry, error) {
	file, err := os.Open(filepath.Join(basePath, "journal", path, "entry.json"))
	if err != nil {
		return Entry{}, err
	}
	var entryDescriptor entryFile
	err = json.NewDecoder(file).Decode(&entryDescriptor)
	if err != nil {
		return Entry{}, err
	}
	journalEntry := Entry{
		Id:          path,
		Comment:     entryDescriptor.Comment,
		Time:        entryDescriptor.Time,
		Laps:        entryDescriptor.Laps,
		LinkedTrack: entryDescriptor.Track,
	}
	date, err := computeDateFromPath(path)
	if err != nil {
		return Entry{}, err
	}
	journalEntry.Date = date
	track, ok := tracks.GetTrack(entryDescriptor.Track)
	if ok {
		journalEntry.Track = &track
	}
	return journalEntry, nil
}

func readListEntry(basePath string, path string) (ListEntry, error) {
	entry, err := ReadJournalEntry(basePath, path)
	if err != nil {
		return ListEntry{}, err
	}
	listEntry := ListEntry{Id: entry.Id, Date: entry.Date}
	if entry.Track != nil {
		listEntry.ParentNames = entry.Track.ParentNames
		listEntry.TrackName = entry.Track.Name
		listEntry.Length = entry.Track.Length * entry.Laps
	}
	return listEntry, nil
}

var dateRegex = regexp.MustCompile("(\\d\\d\\d\\d)-(\\d\\d)-(\\d\\d)")

func CreateEntry(basePath string, date string, trackId string) (ListEntry, error) {
	regexResult := dateRegex.FindStringSubmatch(date)
	if regexResult == nil {
		return ListEntry{}, fmt.Errorf("the date \"%s\" is not a valid date of the format yyyy-mm-dd", date)
	}
	id := filepath.Join(regexResult[1], regexResult[2], regexResult[3])
	journalPath := filepath.Join(basePath, "journal")
	_, existsCheck := os.Stat(filepath.Join(journalPath, id))
	modifier := 0
	for ; existsCheck == nil && modifier < 27; modifier = modifier + 1 {
		_, existsCheck = os.Stat(filepath.Join(journalPath, id+string(rune(modifier+96))))
	}
	if existsCheck == nil {
		return ListEntry{}, fmt.Errorf(
			"all slots for the given date \"%s\" seem to be already taken: %v",
			date,
			existsCheck,
		)
	}
	if modifier > 0 {
		id = id + string(rune(modifier+96))
	}
	err := os.MkdirAll(filepath.Join(basePath, "journal", id), os.ModePerm)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not create directory %s: %v", id, err)
	}
	entryFilePath := filepath.Join(basePath, "journal", id, "entry.json")
	payload, _ := json.Marshal(entryFile{Track: trackId, Laps: 1, Time: "", Comment: ""})
	err = os.WriteFile(entryFilePath, payload, 0644)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not write file \"%s\": %v", entryFilePath, err)
	}
	return readListEntry(basePath, id)
}

func SaveEntry(basePath string, entry Entry) error {
	journalPath := filepath.Join(basePath, "journal", entry.Id, "entry.json")
	if _, err := os.Stat(journalPath); err != nil {
		return fmt.Errorf("could not update entry \"%s\": %v", entry.Id, err)
	}
	payload, _ := json.Marshal(
		entryFile{
			Track:   entry.Track.Id,
			Time:    entry.Time,
			Comment: entry.Comment,
			Laps:    entry.Laps,
		},
	)
	err := os.WriteFile(journalPath, payload, 0644)
	if err != nil {
		return fmt.Errorf("could not update entry \"%s\": %v", entry.Id, err)
	}
	return nil
}

var pathRegex = regexp.MustCompile("(\\d\\d\\d\\d)/(\\d\\d)/(\\d\\d)[a-z]?")

func computeDateFromPath(path string) (string, error) {
	regexResult := pathRegex.FindStringSubmatch(path)
	if regexResult == nil {
		return "", fmt.Errorf("path \"%s\" does not have the correct format dddd/dd/dd/", path)
	}
	return fmt.Sprintf("%s-%s-%s", regexResult[1], regexResult[2], regexResult[3]), nil
}

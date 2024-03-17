package journal

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type ListEntry struct {
	Id          string   `json:"id"`
	Date        string   `json:"date"`
	ParentNames []string `json:"trackParents"`
	TrackName   string   `json:"trackName"`
	Length      int      `json:"length"`
	trackId     string
}

type entryFile struct {
	Track   string `json:"track"`
	Time    string `json:"time"`
	Comment string `json:"comment"`
	Laps    int    `json:"laps"`
}

var directoryRegex = regexp.MustCompile("\\d\\d[a-z]?|\\d\\d\\d\\d")

type Journal struct {
	baseDirectory string
	tracks        map[string]*shared.Track
}

func New(baseDirectory string) (*Journal, error) {
	result := &Journal{baseDirectory: baseDirectory, tracks: make(map[string]*shared.Track)}
	group := sync.WaitGroup{}
	group.Add(1)
	shared.RegisterHandler(
		"tracks initialized", func(data ...any) {
			newTracks, ok := data[0].([]shared.Track)
			if !ok {
				panic(fmt.Errorf("received unexpected type for event \"tracks initialized\": %v", data[0]))
			}
			result.tracks = make(map[string]*shared.Track)
			for _, track := range newTracks {
				t := track
				result.tracks[track.Id] = &t
			}
			group.Done()
		},
	)
	group.Wait()
	entries, err := result.ReadEntries()
	if err == nil {
		for _, entry := range entries {
			shared.Send("journal entry changed", shared.JournalEntry{}, shared.JournalEntry{TrackId: entry.trackId})
		}
	} else {
		return nil, err
	}
	return result, nil
}

func (j *Journal) ReadEntries() ([]ListEntry, error) {
	result := make([]ListEntry, 0, 0)
	journalDirectory := filepath.Join(j.baseDirectory, "journal")
	err := filepath.Walk(
		journalDirectory, func(path string, info os.FileInfo, err error) error {
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
			listEntry, err := j.readListEntry(
				strings.Replace(strings.Replace(path, journalDirectory, "", 1), info.Name(), "", 1),
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
	Track   *shared.Track `json:"track"`
	Comment string        `json:"comment"`
	Time    string        `json:"time"`
	Laps    int           `json:"laps"`
	//needed in FE if the Track file is missing to give a reference what track should be there
	LinkedTrack string `json:"linkedTrack"`
}

func (j *Journal) ReadJournalEntry(path string) (Entry, error) {
	file, err := os.Open(filepath.Join(j.baseDirectory, "journal", path, "entry.json"))
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
	track, ok := j.tracks[entryDescriptor.Track]
	if ok {
		journalEntry.Track = track
	}
	return journalEntry, nil
}

func (j *Journal) readListEntry(path string) (ListEntry, error) {
	entry, err := j.ReadJournalEntry(path)
	if err != nil {
		return ListEntry{}, err
	}
	listEntry := ListEntry{Id: entry.Id, Date: entry.Date}
	if entry.Track != nil {
		listEntry.ParentNames = entry.Track.ParentNames
		listEntry.TrackName = entry.Track.Name
		listEntry.Length = entry.Track.Length * entry.Laps
		listEntry.trackId = entry.LinkedTrack
	}
	return listEntry, nil
}

var dateRegex = regexp.MustCompile("(\\d\\d\\d\\d)-(\\d\\d)-(\\d\\d)")

func (j *Journal) CreateEntry(date string, trackId string) (ListEntry, error) {
	regexResult := dateRegex.FindStringSubmatch(date)
	if regexResult == nil {
		return ListEntry{}, fmt.Errorf("the date \"%s\" is not a valid date of the format yyyy-mm-dd", date)
	}
	id := filepath.Join(regexResult[1], regexResult[2], regexResult[3])
	journalPath := filepath.Join(j.baseDirectory, "journal")
	_, existsCheck := os.Stat(filepath.Join(journalPath, id))
	modifier := 0
	for ; existsCheck == nil && modifier < 27; modifier = modifier + 1 {
		_, existsCheck = os.Stat(filepath.Join(journalPath, id+string(rune(modifier+96))))
	}
	if existsCheck == nil {
		return ListEntry{}, fmt.Errorf(
			"all slots for the given date \"%s\" seem to be already taken: %v", date, existsCheck,
		)
	}
	if modifier > 0 {
		id = id + string(rune(modifier+96))
	}
	err := os.MkdirAll(filepath.Join(j.baseDirectory, "journal", id), os.ModePerm)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not create directory %s: %v", id, err)
	}
	entryFilePath := filepath.Join(j.baseDirectory, "journal", id, "entry.json")
	payload, _ := json.Marshal(entryFile{Track: trackId, Laps: 1, Time: "", Comment: ""})
	err = os.WriteFile(entryFilePath, payload, 0644)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not write file \"%s\": %v", entryFilePath, err)
	}
	return j.readListEntry(id)
}

func (j *Journal) SaveEntry(entry Entry) error {
	journalPath := filepath.Join(j.baseDirectory, "journal", entry.Id, "entry.json")
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

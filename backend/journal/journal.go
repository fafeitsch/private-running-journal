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
	result := &Journal{baseDirectory: filepath.Join(baseDirectory, "journal"), tracks: make(map[string]*shared.Track)}
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
	shared.RegisterHandler(
		"track upserted", func(data ...any) {
			newTrack, ok := data[0].(shared.Track)
			if !ok {
				panic(fmt.Errorf("received unexpected type for event \"track upserted\": %v", data[0]))
			}
			result.tracks[newTrack.Id] = &newTrack
		},
	)
	shared.RegisterHandler(
		"tracks deleted", func(data ...any) {
			deletedTracks, ok := data[0].(map[string]bool)
			if !ok {
				panic(fmt.Errorf("received unexpected type for event \"track deleted\": %v", data[0]))
			}
			for key := range deletedTracks {
				delete(result.tracks, key)
			}
		},
	)
	group.Wait()
	err := os.MkdirAll(result.baseDirectory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create journal directory: %v", err)
	}
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
	err := filepath.Walk(
		j.baseDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			if info.IsDir() || info.Name() != "entry.json" {
				if path == j.baseDirectory || directoryRegex.MatchString(info.Name()) {
					return nil
				}
				log.Printf(
					"skipping directory \"%s\" because it does not have the required format yyyy/mm/dd", path,
				)
				return filepath.SkipDir
			}
			listEntry, err := j.readListEntry(
				strings.Replace(strings.Replace(path, j.baseDirectory+"/", "", 1), "/"+info.Name(), "", 1),
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
	file, err := os.Open(filepath.Join(j.baseDirectory, path, "entry.json"))
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
	path, err := shared.FindFreeFileName(filepath.Join(j.baseDirectory, id))
	id = strings.Replace(path, j.baseDirectory+"/", "", 1)
	if err != nil {
		return ListEntry{}, err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not create directory %s: %v", id, err)
	}
	entryFilePath := filepath.Join(j.baseDirectory, id, "entry.json")
	payload, _ := json.Marshal(entryFile{Track: trackId, Laps: 1, Time: "", Comment: ""})
	err = os.WriteFile(entryFilePath, payload, 0644)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not write file \"%s\": %v", entryFilePath, err)
	}
	shared.Send("journal entry changed", shared.JournalEntry{}, shared.JournalEntry{TrackId: trackId})
	return j.readListEntry(id)
}

func (j *Journal) SaveEntry(entry Entry) error {
	journalPath := filepath.Join(j.baseDirectory, entry.Id, "entry.json")
	oldEntry, err := j.readListEntry(entry.Id)
	if err != nil {
		return fmt.Errorf("could not read old entry \"%s\": %v", entry.Id, err)
	}
	payload, _ := json.Marshal(
		entryFile{
			Track:   entry.Track.Id,
			Time:    entry.Time,
			Comment: entry.Comment,
			Laps:    entry.Laps,
		},
	)
	err = os.WriteFile(journalPath, payload, 0644)
	if err != nil {
		return fmt.Errorf("could not update entry \"%s\": %v", entry.Id, err)
	}
	shared.Send(
		"journal entry changed",
		shared.JournalEntry{TrackId: oldEntry.trackId},
		shared.JournalEntry{TrackId: entry.Track.Id},
	)
	return nil
}

func (j *Journal) DeleteEntry(key string) error {
	entry, err := j.readListEntry(key)
	if err != nil {
		return nil
	}
	journalPath := filepath.Join(j.baseDirectory, key)
	shared.Send("journal entry changed", shared.JournalEntry{TrackId: entry.trackId}, shared.JournalEntry{})
	return os.RemoveAll(journalPath)
}

var pathRegex = regexp.MustCompile("(\\d\\d\\d\\d)/(\\d\\d)/(\\d\\d)[a-z]?")

func computeDateFromPath(path string) (string, error) {
	regexResult := pathRegex.FindStringSubmatch(path)
	if regexResult == nil {
		return "", fmt.Errorf("path \"%s\" does not have the correct format dddd/dd/dd/", path)
	}
	return fmt.Sprintf("%s-%s-%s", regexResult[1], regexResult[2], regexResult[3]), nil
}

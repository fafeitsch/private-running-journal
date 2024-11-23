package journal

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ListEntry struct {
	Id        string `json:"id"`
	Date      string `json:"date"`
	TrackName string `json:"trackName"`
	Length    int    `json:"length"`
	trackId   string
}

type entryFile struct {
	Track        string `json:"track"`
	Time         string `json:"time"`
	Comment      string `json:"comment"`
	Laps         int    `json:"laps"`
	CustomLength int    `json:"customLength,omitempty"`
}

var directoryRegex = regexp.MustCompile("\\d\\d[a-z]?|\\d\\d\\d\\d")

type Journal struct {
	baseDirectory string
	listEntries   map[string]*ListEntry
	loadService   *filebased.Service
}

func New(baseDirectory string, loadService *filebased.Service) (
	*Journal, error,
) {
	result := &Journal{
		baseDirectory: filepath.Join(baseDirectory, "journal"),
		listEntries:   make(map[string]*ListEntry),
		loadService:   loadService,
	}
	shared.RegisterHandler(
		"track deleted", func(data ...any) {
			deletedTrack, ok := data[0].(string)
			if !ok {
				panic(fmt.Errorf("received unexpected type for event \"track deleted\": %v", data[0]))
			}
			for key, entry := range result.listEntries {
				if entry.trackId == deletedTrack {
					result.listEntries[key].TrackName = ""
					result.listEntries[key].Length = 0
				}
			}
		},
	)
	shared.RegisterHandler(
		"track moved", func(data ...any) {
			oldId, ok1 := data[0].(string)
			track, ok2 := data[1].(shared.Track)
			if !ok1 || !ok2 {
				panic(
					fmt.Errorf(
						"expected the old id and the new track as the two event params, but was %v and %v",
						data[0],
						data[1],
					),
				)
			}
			for _, listEntry := range result.listEntries {
				entry, err := result.ReadJournalEntry(listEntry.Id)
				if entry.LinkedTrack == oldId {
					entry.LinkedTrack = track.Id
					entry.Track = &track
					err = result.SaveEntry(entry)
					if err != nil {
						log.Printf("could not update track id of %s: %v", listEntry.Id, err)
					}
				}
			}
		},
	)
	entries, err := result.readEntries()
	if err == nil {
		for _, entry := range entries {
			listEntry := entry
			result.listEntries[entry.Id] = &listEntry
			shared.Send("journal entry loaded", shared.JournalEntry{TrackId: entry.trackId})
		}
	} else {
		return nil, err
	}
	return result, nil
}

func (j *Journal) ReadAllEntries() ([]shared.JournalEntry, error) {
	entries, err := j.readEntries()
	if err != nil {
		return nil, err
	}
	result := make([]shared.JournalEntry, 0)
	for _, entry := range entries {
		result = append(result, shared.JournalEntry{TrackId: entry.trackId, Id: entry.Id})
	}
	return result, nil
}

func (j *Journal) readEntries() ([]ListEntry, error) {
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
	Id           string        `json:"id"`
	Date         string        `json:"date"`
	Track        *shared.Track `json:"track"`
	Comment      string        `json:"comment"`
	Time         string        `json:"time"`
	Laps         int           `json:"laps"`
	CustomLength int           `json:"customLength,omitempty"`
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
		Id:           path,
		Comment:      entryDescriptor.Comment,
		Time:         entryDescriptor.Time,
		Laps:         entryDescriptor.Laps,
		LinkedTrack:  entryDescriptor.Track,
		CustomLength: entryDescriptor.CustomLength,
	}
	date, err := computeDateFromPath(path)
	if err != nil {
		return Entry{}, err
	}
	journalEntry.Date = date
	track, err := j.loadService.ReadTrack(entryDescriptor.Track)
	if err == nil {
		journalEntry.Track = &track
	} else {
		log.Printf("could not load track: %v", err)
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
		listEntry.TrackName = entry.Track.Name
		listEntry.Length = entry.Track.Waypoints.Length() * entry.Laps
	}
	listEntry.trackId = entry.LinkedTrack
	if entry.CustomLength != 0 {
		listEntry.Length = entry.CustomLength
	}
	return listEntry, nil
}

var dateRegex = regexp.MustCompile("(\\d\\d\\d\\d)-(\\d\\d)-(\\d\\d)")

func (j *Journal) CreateEntry(entry Entry) (ListEntry, error) {
	regexResult := dateRegex.FindStringSubmatch(entry.Date)
	if regexResult == nil {
		return ListEntry{}, fmt.Errorf("the date \"%s\" is not a valid date of the format yyyy-mm-dd", entry.Date)
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
	payload, _ := json.Marshal(
		entryFile{
			Track:        entry.Track.Id,
			Laps:         entry.Laps,
			Time:         entry.Time,
			Comment:      entry.Comment,
			CustomLength: entry.CustomLength,
		},
	)
	err = os.WriteFile(entryFilePath, payload, 0644)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not write file \"%s\": %v", entryFilePath, err)
	}
	shared.Send("journal entry changed", shared.JournalEntry{}, shared.JournalEntry{TrackId: entry.Track.Id, Id: id})
	listEntry, err := j.readListEntry(id)
	if err != nil {
		return ListEntry{}, fmt.Errorf("could not read file \"%s\": %v", entryFilePath, err)
	}
	j.listEntries[id] = &listEntry
	return *j.listEntries[id], nil
}

func (j *Journal) SaveEntry(entry Entry) error {
	journalPath := filepath.Join(j.baseDirectory, entry.Id, "entry.json")
	oldEntry, err := j.readListEntry(entry.Id)
	if err != nil {
		return fmt.Errorf("could not read old entry \"%s\": %v", entry.Id, err)
	}
	payload, _ := json.Marshal(
		entryFile{
			Track:        entry.Track.Id,
			Time:         entry.Time,
			Comment:      entry.Comment,
			Laps:         entry.Laps,
			CustomLength: entry.CustomLength,
		},
	)
	err = os.WriteFile(journalPath, payload, 0644)
	if err != nil {
		return fmt.Errorf("could not update entry \"%s\": %v", entry.Id, err)
	}
	shared.Send(
		"journal entry changed",
		shared.JournalEntry{TrackId: oldEntry.trackId, Id: entry.Id},
		shared.JournalEntry{TrackId: entry.Track.Id, Id: entry.Id},
	)
	listEntry, err := j.readListEntry(entry.Id)
	if err != nil {
		return fmt.Errorf("could not read entry file \"%s\": %v", entry.Id, err)
	}
	j.listEntries[entry.Id] = &listEntry
	return nil
}

func (j *Journal) DeleteEntry(key string) error {
	entry, err := j.readListEntry(key)
	if err != nil {
		return nil
	}
	journalPath := filepath.Join(j.baseDirectory, key)
	err = os.RemoveAll(journalPath)
	if err == nil {
		shared.Send("journal entry changed", shared.JournalEntry{TrackId: entry.trackId}, shared.JournalEntry{})
		delete(j.listEntries, key)
	}
	return err
}

var pathRegex = regexp.MustCompile("(\\d\\d\\d\\d)/(\\d\\d)/(\\d\\d)[a-z]?")

func computeDateFromPath(path string) (string, error) {
	regexResult := pathRegex.FindStringSubmatch(path)
	if regexResult == nil {
		return "", fmt.Errorf("path \"%s\" does not have the correct format dddd/dd/dd/", path)
	}
	return fmt.Sprintf("%s-%s-%s", regexResult[1], regexResult[2], regexResult[3]), nil
}

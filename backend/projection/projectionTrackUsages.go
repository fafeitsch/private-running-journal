package projection

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/journal"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
)

type TrackUsagesProjector struct {
	Journal   *journal.Journal
	retriever Retriever
	rebuilder Rebuilder
}

func (t *TrackUsagesProjector) Bootstrap(retriever Retriever, rebuilder Rebuilder) {
	t.retriever = retriever
	t.rebuilder = rebuilder
	shared.RegisterHandler(
		"journal entry changed", func(data ...any) {
			old := data[0].(shared.JournalEntry)
			nevv := data[1].(shared.JournalEntry)
			message, err := t.retriever()
			if err != nil {
				return
			}
			var result map[string][]string
			err = json.Unmarshal(message, &result)
			if err != nil {
				return
			}
			log.Printf("old usages", result)
			if old.TrackId == nevv.TrackId {
				return
			}
			log.Printf("new track unqueal old track")
			oldUsages, ok := result[old.TrackId]
			if ok {
				updatedOldUsages := make([]string, 0)
				filtered := false
				for _, usage := range oldUsages {
					fmt.Printf("usage: %v, id: %v", usage, nevv.Id)
					if usage == nevv.Id && !filtered {
						filtered = true
						continue
					}
					updatedOldUsages = append(updatedOldUsages, usage)
				}
				result[old.TrackId] = updatedOldUsages
			}
			_, ok = result[nevv.TrackId]
			if ok {
				result[nevv.TrackId] = append(result[nevv.TrackId], nevv.Id)
			}
			payload, err := json.Marshal(result)
			_ = rebuilder(payload)
		},
	)
}

func (t *TrackUsagesProjector) BuildProjection() (json.RawMessage, error) {
	result := make(map[string][]string)
	entries, err := t.Journal.ReadAllEntries()
	if err != nil {
		return nil, fmt.Errorf("could not read journal entries: %v", err)
	}
	for _, entry := range entries {
		if _, ok := result[entry.TrackId]; !ok {
			result[entry.TrackId] = make([]string, 0)
		}
		result[entry.TrackId] = append(result[entry.TrackId], entry.Id)
	}
	message, _ := json.Marshal(result)
	return message, nil
}

func (t *TrackUsagesProjector) ProjectionName() string {
	return "trackUsages"
}

func (t *TrackUsagesProjector) GetUsages(trackId string) ([]string, error) {
	log.Printf("Formatter %v", t)
	message, err := t.retriever()
	if err != nil {
		return nil, err
	}
	var result map[string][]string
	err = json.Unmarshal(message, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal usages: %v", err)
	}
	return result[trackId], nil
}

package projection

import (
	"encoding/json"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"slices"
	"sync"
)

type TrackUsagesMap map[string][]string

type TrackUsages struct {
	sync.RWMutex
	content TrackUsagesMap
}

func (t *TrackUsages) ProjectionName() string {
	return "trackUsages"
}

func (t *TrackUsages) Init(message json.RawMessage, writer func()) {
	if message != nil {
		_ = json.Unmarshal(message, &t.content)
	} else {
		t.content = make(map[string][]string)
	}
	shared.Listen(
		shared.JournalEntryUpsertedEvent{}, func(event shared.JournalEntryUpsertedEvent) {
			t.Lock()
			old := event.OldTrackId
			nevv := event.TrackId
			if old == nevv {
				return
			}
			if _, ok := t.content[old]; ok {
				filtered := false
				t.content[old] = slices.DeleteFunc(
					t.content[old], func(s string) bool {
						result := s == event.Id && !filtered
						filtered = s == event.Id
						return result
					},
				)
			}

			if _, ok := t.content[nevv]; ok {
				t.content[nevv] = append(t.content[nevv], event.Id)
			}
			t.Unlock()
			writer()
		},
	)
	shared.Listen(
		shared.TrackDeletedEvent{}, func(event shared.TrackDeletedEvent) {
			t.Lock()
			delete(t.content, event.Id)
			t.Unlock()
			writer()
		},
	)
}

func (t *TrackUsages) AddTrack(track shared.Track) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.content[track.Id]; !ok {
		t.content[track.Id] = make([]string, 0)
	}
}

func (t *TrackUsages) AddJournalEntry(entry shared.JournalEntry) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.content[entry.TrackId]; !ok {
		t.content[entry.TrackId] = make([]string, 0)
	}
	t.content[entry.TrackId] = append(t.content[entry.TrackId], entry.Id)
}

func (t *TrackUsages) GetData() any {
	t.RLock()
	defer t.RUnlock()
	return t.content
}

func (t *TrackUsages) Get() TrackUsagesMap {
	t.RLock()
	defer t.RUnlock()
	return t.content
}

func (t *TrackUsages) GetUsages(trackId string) ([]string, error) {
	t.RLock()
	defer t.RUnlock()
	usages, ok := t.content[trackId]
	if !ok {
		usages = make([]string, 0)
	}
	return usages, nil
}

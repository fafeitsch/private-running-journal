package projection

import (
	"encoding/json"
	"github.com/fafeitsch/private-running-journal/backend/shared"
)

type TrackIdMap map[string][]string

type TrackLookup struct {
	projection TrackIdMap
}

func (t *TrackLookup) Init(message json.RawMessage, writer func()) {
	if message != nil {
		_ = json.Unmarshal(message, &t.projection)
	} else {
		t.projection = TrackIdMap{}
	}
	shared.Listen(
		shared.TrackUpsertedEvent{}, func(newTrack shared.TrackUpsertedEvent) {
			t.projection[newTrack.Id] = newTrack.Parents
			writer()
		},
	)
	shared.Listen(
		shared.TrackDeletedEvent{}, func(event shared.TrackDeletedEvent) {
			delete(t.projection, event.Id)
			writer()
		},
	)
}

func (t *TrackLookup) AddTrack(track shared.Track) {
	t.projection[track.Id] = track.Parents
}

func (t *TrackLookup) AddJournalEntry(entry shared.JournalEntry) {
}

func (t *TrackLookup) GetData() any {
	return t.projection
}

func (t *TrackLookup) Get() TrackIdMap {
	return t.projection
}

func (t *TrackLookup) ProjectionName() string {
	return "trackIdMap"
}

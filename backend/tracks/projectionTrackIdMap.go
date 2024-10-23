package tracks

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
)

type TrackIdMapProjection struct {
	*Tracks
	retriever projection.Retriever
	rebuilder projection.Rebuilder
}

func (t *TrackIdMapProjection) ProjectionName() string {
	return "trackIdMap"
}

func (t *TrackIdMapProjection) Bootstrap(retriever projection.Retriever, rebuilder projection.Rebuilder) {
	t.retriever = retriever
	shared.Listen(
		shared.TrackUpsertedEvent{}, func(newTrack shared.TrackUpsertedEvent) {
			current, err := t.retriever()
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
			var trackMap map[string][]string
			err = json.Unmarshal(current, &trackMap)
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
			trackMap[newTrack.Id] = newTrack.Parents
			trackMap[newTrack.Id] = append(trackMap[newTrack.Id], newTrack.Name)
			message, _ := json.Marshal(trackMap)
			err = rebuilder(message)
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
		},
	)
}

func (t *TrackIdMapProjection) BuildProjection() (json.RawMessage, error) {
	cache := make(map[string][]string)
	err := t.walkTracksDirectory(
		func(track Track) {
			cache[track.Id] = track.Hierarchy
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not walk directory %s: %v", t.basePath, err)
	}
	return json.Marshal(cache)
}

func (t *TrackIdMapProjection) GetTrackLocation(id string) ([]string, error) {
	var result map[string][]string
	message, err := t.retriever()
	if err != nil {
		return nil, fmt.Errorf("could not open track list: %v", err)
	}
	err = json.Unmarshal(message, &result)
	return result[id], err
}

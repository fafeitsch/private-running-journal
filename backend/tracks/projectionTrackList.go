package tracks

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"golang.org/x/exp/slices"
	"log"
)

type TrackListEntry struct {
	Name   string
	Id     string
	Length int
}

type trackListProjector struct {
	*Tracks
	retriever projection.Retriever
	rebuilder projection.Rebuilder
}

func (t *trackListProjector) ProjectionName() string {
	return "trackList"
}

func (t *trackListProjector) Bootstrap(retriever projection.Retriever, rebuilder projection.Rebuilder) {
	t.retriever = retriever
	shared.RegisterHandler(
		"track upserted", func(data ...any) {
			newTrack, ok := data[0].(shared.Track)
			if !ok {
				panic(fmt.Errorf("received unexpected type for event \"track upserted\": %v", data[0]))
			}
			current, err := t.retriever()
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
			var currentList []TrackListEntry
			err = json.Unmarshal(current, &currentList)
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
			index := slices.IndexFunc(
				currentList, func(e TrackListEntry) bool {
					return e.Id == newTrack.Id
				},
			)
			if index > -1 {
				currentList[index].Name = newTrack.Name
				currentList[index].Length = newTrack.Length
			} else {
				currentList = append(currentList, TrackListEntry{newTrack.Name, newTrack.Id, newTrack.Length})
			}
			message, _ := json.Marshal(currentList)
			err = rebuilder(message)
			if err != nil {
				log.Printf("could not update trackList projection after upserting: %v", err)
			}
		},
	)
}

func (t *trackListProjector) BuildProjection() (json.RawMessage, error) {
	cache := make(map[string]Track)
	err := t.walkTracksDirectory(
		func(track Track) {
			cache[track.Id] = track
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not walk directory %s: %v", t.basePath, err)
	}
	result := make([]TrackListEntry, 0)
	for _, track := range cache {
		result = append(
			result, TrackListEntry{
				Name:   track.Name,
				Id:     track.Id,
				Length: track.Length,
			},
		)
	}
	return json.Marshal(result)
}

func (t *trackListProjector) loadTrackList() ([]TrackListEntry, error) {
	var result = make([]TrackListEntry, 0)
	message, err := t.retriever()
	if err != nil {
		return nil, fmt.Errorf("could not open track list: %v", err)
	}
	err = json.Unmarshal(message, &result)
	return result, err
}

package tracks

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"golang.org/x/exp/slices"
	"log"
)

type TrackTreeNode struct {
	Name   string           `json:"name"`
	Tracks []TrackTreeEntry `json:"tracks"`
	Nodes  []TrackTreeNode  `json:"nodes"`
}

type TrackTreeEntry struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type trackTreeProjector struct {
	*Tracks
	retriever projection.Retriever
	rebuilder projection.Rebuilder
}

func (t *trackTreeProjector) ProjectionName() string {
	return "trackTree"
}

func (t *trackTreeProjector) Bootstrap(retriever projection.Retriever, rebuilder projection.Rebuilder) {
	t.retriever = retriever
	t.rebuilder = rebuilder
	shared.RegisterHandler(
		"track upserted", func(data ...any) {
			message, err := t.BuildProjection()
			if err != nil {
				log.Printf("could not update trackTree projection after upserting: %v", err)
			}
			err = rebuilder(message)
			if err != nil {
				log.Printf("could not update trackTree projection after upserting: %v", err)
			}
		},
	)
}

func (t *trackTreeProjector) BuildProjection() (json.RawMessage, error) {
	tree := TrackTreeNode{Tracks: make([]TrackTreeEntry, 0), Nodes: make([]TrackTreeNode, 0)}
	err := t.walkTracksDirectory(
		func(track Track) {
			hierarchy := track.Hierarchy
			node := &tree
			for len(hierarchy) > 0 {
				index := slices.IndexFunc(
					node.Nodes, func(node TrackTreeNode) bool {
						return node.Name == hierarchy[0]
					},
				)
				if index == -1 {
					newNode := TrackTreeNode{
						Name:   hierarchy[0],
						Tracks: make([]TrackTreeEntry, 0),
						Nodes:  make([]TrackTreeNode, 0),
					}
					node.Nodes = append(node.Nodes, newNode)
					node = &node.Nodes[len(node.Nodes)-1]
				} else {
					node = &node.Nodes[index]
				}
				hierarchy = hierarchy[1:]
			}
			node.Tracks = append(node.Tracks, TrackTreeEntry{track.Name, track.Id})
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not walk directory %s: %v", t.basePath, err)
	}
	log.Printf("Tree %v", tree)
	return json.Marshal(tree)
}

func (t *trackTreeProjector) loadTrackTree() (TrackTreeNode, error) {
	var result TrackTreeNode
	message, err := t.retriever()
	if err != nil {
		return TrackTreeNode{}, fmt.Errorf("could not open track tree: %v", err)
	}
	err = json.Unmarshal(message, &result)
	return result, err
}

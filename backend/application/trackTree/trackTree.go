package trackTree

import (
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"strings"
)

type TrackTree struct {
	trackTreeProjector *projection.TrackTree
}

func New(trackTreeProjector *projection.TrackTree) *TrackTree {
	return &TrackTree{
		trackTreeProjector: trackTreeProjector,
	}
}

func (t *TrackTree) GetTrackTree(name string) projection.TrackTreeNode {
	if name == "" {
		return t.trackTreeProjector.Get()
	}

	originalTree := t.trackTreeProjector.Get()

	filteredTree := t.filterNode(&originalTree, strings.ToLower(name))

	if filteredTree == nil {
		return projection.TrackTreeNode{
			Name:   originalTree.Name,
			Tracks: make([]projection.TrackTreeEntry, 0),
			Nodes:  make([]*projection.TrackTreeNode, 0),
		}
	}

	return *filteredTree
}

func (t *TrackTree) filterNode(original *projection.TrackTreeNode, name string) *projection.TrackTreeNode {
	filteredTracks := make([]projection.TrackTreeEntry, 0)
	filteredNodes := make([]*projection.TrackTreeNode, 0)

	for _, track := range original.Tracks {
		if strings.Contains(strings.ToLower(track.Name), name) {
			filteredTracks = append(filteredTracks, track)
		}
	}

	for _, node := range original.Nodes {
		filteredNode := t.filterNode(node, name)
		if filteredNode != nil && (len(filteredNode.Tracks) > 0 || len(filteredNode.Nodes) > 0) {
			filteredNodes = append(filteredNodes, filteredNode)
		}
	}

	if len(filteredTracks) == 0 && len(filteredNodes) == 0 {
		return nil
	}

	return &projection.TrackTreeNode{
		Name:   original.Name,
		Tracks: filteredTracks,
		Nodes:  filteredNodes,
	}
}

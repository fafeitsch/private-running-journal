package projection

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"slices"
)

type TrackTreeNode struct {
	Name   string           `json:"name"`
	Tracks []TrackTreeEntry `json:"tracks"`
	Nodes  []*TrackTreeNode `json:"nodes"`
}

type TrackTreeEntry struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Length int    `json:"length"`
}

type TrackTree struct {
	fileService *filebased.Service
	tree        *TrackTreeNode
}

func (t *TrackTree) ProjectionName() string {
	return "trackTree"
}

func (t *TrackTree) Init(message json.RawMessage, writer func()) {
	if message == nil {
		t.tree = &TrackTreeNode{Tracks: make([]TrackTreeEntry, 0), Nodes: make([]*TrackTreeNode, 0)}
	} else {
		_ = json.Unmarshal(message, &t.tree)
	}
	shared.Listen(
		shared.TrackUpsertedEvent{}, func(k shared.TrackUpsertedEvent) {
			t.handleUpsertEvent(k)
			writer()
		},
	)
	shared.Listen(
		shared.TrackDeletedEvent{}, func(k shared.TrackDeletedEvent) {
			t.handleDeleteEvent(t.tree, k.Id)
			writer()
		},
	)
}

func (t *TrackTree) AddTrack(track shared.Track) {
	hierarchy := track.Parents
	node := t.tree
	for len(hierarchy) > 0 {
		index := slices.IndexFunc(
			node.Nodes, func(node *TrackTreeNode) bool {
				return node.Name == hierarchy[0]
			},
		)
		if index == -1 {
			newNode := &TrackTreeNode{
				Name:   hierarchy[0],
				Tracks: make([]TrackTreeEntry, 0),
				Nodes:  make([]*TrackTreeNode, 0),
			}
			node.Nodes = append(node.Nodes, newNode)
			node = node.Nodes[len(node.Nodes)-1]
		} else {
			node = node.Nodes[index]
		}
		hierarchy = hierarchy[1:]
	}
	node.Tracks = append(node.Tracks, TrackTreeEntry{Name: track.Name, Id: track.Id, Length: track.Waypoints.Length()})
}

func (t *TrackTree) AddJournalEntry(entry shared.JournalEntry) {}

func (t *TrackTree) Get() TrackTreeNode {
	return *t.tree
}

func (t *TrackTree) GetData() any {
	return t.tree
}

func (t *TrackTree) handleDeleteEvent(root *TrackTreeNode, id string) {
	root.Tracks = slices.DeleteFunc(
		root.Tracks, func(track TrackTreeEntry) bool {
			return track.Id == id
		},
	)
	for i := range root.Nodes {
		t.handleDeleteEvent(root.Nodes[i], id)
	}
}

func (t *TrackTree) handleUpsertEvent(track shared.TrackUpsertedEvent) {
	t.handleDeleteEvent(t.tree, track.Id)
	queue := track.Parents[:len(track.Parents)-1]
	node := t.tree
outer:
	for len(queue) > 0 {
		head := queue[0]
		for _, n := range node.Nodes {
			if n.Name == head {
				queue = queue[1:]
				node = n
				continue outer
			}
		}
		node.Nodes = append(
			node.Nodes,
			&TrackTreeNode{Name: head, Tracks: make([]TrackTreeEntry, 0), Nodes: make([]*TrackTreeNode, 0)},
		)
	}
	fmt.Printf("parents: %v, node: %v", track.Parents, node.Name)
	node.Tracks = append(node.Tracks, TrackTreeEntry{Name: track.Name, Id: track.Id, Length: track.Waypoints.Length()})
}

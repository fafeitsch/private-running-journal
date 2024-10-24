package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/twpayne/go-gpx"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Track struct {
	Id        string   `json:"id"`
	Length    int      `json:"length"`
	Name      string   `json:"name"`
	Usages    int      `json:"usages"`
	Hierarchy []string `json:"hierarchy"`
}

type Tracks struct {
	basePath             string
	loadService          *filebased.Service
	Projectors           []projection.Projector
	TrackIdMapProjector  *TrackIdMapProjection
	trackTreeProjector   *trackTreeProjector
	TrackUsagesProjector *projection.TrackUsagesProjector
}

func New(baseDir string, loadService *filebased.Service) *Tracks {
	result := Tracks{basePath: filepath.Join(baseDir, "tracks"), loadService: loadService}
	result.TrackIdMapProjector = &TrackIdMapProjection{Tracks: &result}
	result.trackTreeProjector = &trackTreeProjector{Tracks: &result}
	result.Projectors = []projection.Projector{result.trackTreeProjector, result.TrackIdMapProjector}
	return &result
}

func (t *Tracks) readTrack(path string, relativePath string) (Track, error) {
	descriptorPath := filepath.Join(path, "info.json")
	var baseDescriptor trackDescriptor
	fileContent, err := os.Open(descriptorPath)
	if err != nil {
		return Track{}, err
	}
	err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
	if err != nil {
		return Track{}, err
	}

	length := 0
	gpxPath := filepath.Join(path, "track.gpx")
	if _, err := os.Stat(gpxPath); err == nil {
		_, length, err = readGpx(gpxPath)
		if err != nil {
			return Track{}, err
		}
	}

	hierarchy := strings.Split(relativePath, string(os.PathSeparator))
	return Track{
		Id:        baseDescriptor.Id,
		Length:    length,
		Name:      baseDescriptor.Name,
		Hierarchy: hierarchy,
	}, nil
}

func (t *Tracks) TrackTree() (TrackTreeNode, error) {
	return t.trackTreeProjector.loadTrackTree()
}

type SaveTrack struct {
	Id        string        `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Waypoints []Coordinates `json:"waypoints,omitempty"`
}

func (t *Tracks) GetTrack(id string) (Track, error) {
	result, err := t.readTrack(path.Join(t.basePath, id), id)
	if err != nil {
		return Track{}, err
	}
	usages, err := t.TrackUsagesProjector.GetUsages(id)
	result.Usages = len(usages)
	return result, nil
}

type CreateTrack struct {
	*SaveTrack
	Parent string `json:"parent"`
}

type trackDescriptor struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func readGpx(path string) ([]Coordinates, int, error) {
	gpxFileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("ocould not read gpx track %s: %v", path, err)
	}
	tracks, err := gpx.Read(bytes.NewReader(gpxFileContent))
	if err != nil {
		return nil, 0, fmt.Errorf("could not parse gpx %s: %v", path, err)
	}
	if len(tracks.Trk) != 1 {
		return nil, 0, fmt.Errorf("%s must contain one track only, but contains %d tracks", path, len(tracks.Trk))
	}
	length := 0.0
	coordinates := make([]Coordinates, 0, 0)
	for _, segment := range tracks.Trk[0].TrkSeg {
		length = length + distanceGpx(segment.TrkPt)
		for _, trkPt := range segment.TrkPt {
			coordinates = append(coordinates, Coordinates{Longitude: trkPt.Lon, Latitude: trkPt.Lat})
		}
	}
	return coordinates, int(1000 * length), nil
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (t *Tracks) walkTracksDirectory(consumer func(track Track)) error {
	return filepath.WalkDir(
		t.basePath, func(path string, info os.DirEntry, err error) error {
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			if info.IsDir() || info.Name() != "info.json" {
				return nil
			}
			parent := strings.Replace(path, "/"+info.Name(), "", 1)
			relativePath := strings.Replace(parent, t.basePath+"/", "", 1)
			track, err := t.readTrack(parent, relativePath)
			if track.Id == "" {
				track.Id = relativePath
			}
			if err != nil {
				log.Printf("skipping track \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			consumer(track)
			return nil
		},
	)
}

package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/twpayne/go-gpx"
	"io"
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
		Id:        relativePath,
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

func (t *Tracks) DeleteTrack(id string) error {
	err := os.RemoveAll(filepath.Join(t.basePath, id))
	t.deleteEmptyDirectories(id)
	shared.Send("track deleted", id)
	return err
}

func (t *Tracks) MoveTrack(id string, newPath string) (*Track, error) {
	_, err := t.readTrack(path.Join(t.basePath, id), id)
	if err != nil {
		return nil, fmt.Errorf("could not find track with id %s", id)
	}
	lastSegment := id[strings.LastIndex(id, string(filepath.Separator))+1:]
	newDirPath := filepath.Join(t.basePath, newPath)
	err = os.MkdirAll(newDirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create directories: %v", err)
	}
	err = os.Rename(filepath.Join(t.basePath, id), filepath.Join(newDirPath, lastSegment))
	if err != nil {
		return nil, fmt.Errorf("could not move track: %v", err)
	}
	movedTrack, err := t.readTrack(filepath.Join(newDirPath, lastSegment), filepath.Join(newPath, lastSegment))
	t.deleteEmptyDirectories(id)
	shared.Send(
		"track moved", id, shared.Track{
			Id:     movedTrack.Id,
			Length: movedTrack.Length,
			Name:   movedTrack.Name,
		},
	)
	return &movedTrack, err
}

func (t *Tracks) deleteEmptyDirectories(id string) {
	parts := strings.Split(id, string(filepath.Separator))
	if len(parts) == 0 {
		return
	}
	counter := 1
	directory := filepath.Join(parts[:len(parts)-counter]...)
	file, err := os.Open(filepath.Join(t.basePath, directory))
	if err != nil {
		log.Printf("could not check whether directory is empty: %v", err)
		return
	}
	_, err = file.Readdirnames(1)
	_ = file.Close()
	for err == io.EOF && counter < len(parts) {
		_ = file.Close()
		log.Printf("deleting %v", directory)
		_ = os.RemoveAll(filepath.Join(t.basePath, directory))
		counter = counter + 1
		directory = filepath.Join(parts[:len(parts)-counter]...)
		file, err = os.Open(filepath.Join(t.basePath, directory))
		if err != nil {
			log.Printf("could not check whether directory is empty: %v", err)
			return
		}
		_, err = file.Readdirnames(1)
		log.Printf("%s %v", directory, err)
	}
	_ = file.Close()
}

type trackDescriptor struct {
	Name string `json:"name"`
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
			if err != nil {
				log.Printf("skipping track \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			consumer(track)
			return nil
		},
	)
}

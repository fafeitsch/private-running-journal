package tracks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/twpayne/go-geom"
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
	trackListProjector   *trackListProjector
	trackTreeProjector   *trackTreeProjector
	TrackUsagesProjector *projection.TrackUsagesProjector
}

func New(baseDir string, loadService *filebased.Service) *Tracks {
	result := Tracks{basePath: filepath.Join(baseDir, "tracks"), loadService: loadService}
	result.trackListProjector = &trackListProjector{Tracks: &result}
	result.trackTreeProjector = &trackTreeProjector{Tracks: &result}
	result.Projectors = []projection.Projector{result.trackTreeProjector, result.trackListProjector}
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
		Hierarchy: hierarchy[:len(hierarchy)-1],
	}, nil
}

func (t *Tracks) Tracks() ([]TrackListEntry, error) {
	return t.trackListProjector.loadTrackList()
}

func (t *Tracks) TrackTree() (TrackTreeNode, error) {
	return t.trackTreeProjector.loadTrackTree()
}

type SaveTrack struct {
	Id        string        `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Waypoints []Coordinates `json:"waypoints,omitempty"`
}

func (t *Tracks) SaveTrack(track SaveTrack) (*Track, error) {
	trackDirectory := path.Join(t.basePath, track.Id)
	stat, err := os.Stat(trackDirectory)
	if err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("derived track directory \"%s\" does not seems to exist: %v", trackDirectory, err)
	}
	existing, err := t.readTrack(path.Join(t.basePath, track.Id), track.Id)
	if err != nil {
		return nil, fmt.Errorf("the track with id \"%s\" does not seem to exist yet", track.Id)
	}
	existing.Name = track.Name
	existing.Length = int(1000 * distance(track.Waypoints))
	infoFile := path.Join(trackDirectory, "info.json")
	infoPayload, _ := json.Marshal(trackDescriptor{Name: existing.Name})
	err = os.WriteFile(infoFile, infoPayload, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not save base information: %v", err)
	}
	err = writeGpxFile(track, trackDirectory)
	if err != nil {
		return nil, fmt.Errorf("could not save gpx file: %v", err)
	}
	existingTrack, err := t.readTrack(trackDirectory, track.Id)
	if err != nil {
		return nil, fmt.Errorf("could not read track: %v", err)
	}
	shared.Send(
		"track upserted", shared.Track{
			Id:     existingTrack.Id,
			Length: existingTrack.Length,
			Name:   existingTrack.Name,
		},
	)
	return &existingTrack, err
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

func writeGpxFile(track SaveTrack, trackDirectory string) error {
	coords := make([]geom.Coord, 0)
	for _, coordinate := range track.Waypoints {
		coords = append(coords, []float64{coordinate.Longitude, coordinate.Latitude})
	}
	linestring, _ := geom.NewLineString(geom.XY).SetCoords(coords)
	segment := gpx.NewTrkSegType(linestring)
	trackSegment := &gpx.TrkType{TrkSeg: []*gpx.TrkSegType{segment}}
	gpxPayload := gpx.GPX{Trk: []*gpx.TrkType{trackSegment}}
	writer := bytes.Buffer{}
	_ = gpxPayload.WriteIndent(bufio.NewWriter(&writer), "  ", "  ")
	return os.WriteFile(filepath.Join(trackDirectory, "track.gpx"), writer.Bytes(), 0644)
}

type CreateTrack struct {
	*SaveTrack
	Parent string `json:"parent"`
}

func (t *Tracks) CreateTrack(track CreateTrack) (*Track, error) {
	parentId := track.Parent
	trackPath, err := shared.FindFreeFileName(filepath.Join(t.basePath, parentId, strings.ToLower(track.Name)))
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(trackPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	entryFilePath := filepath.Join(trackPath, "info.json")
	payload, _ := json.Marshal(trackDescriptor{Name: track.Name})
	err = os.WriteFile(entryFilePath, payload, 0644)
	if err != nil {
		return nil, err
	}
	err = writeGpxFile(*track.SaveTrack, trackPath)
	if err != nil {
		return nil, err
	}
	id := strings.Replace(trackPath, t.basePath+"/", "", 1)
	newTrack, err := t.readTrack(trackPath, id)
	if err != nil {
		return nil, err
	}
	shared.Send(
		"track upserted", shared.Track{
			Id:     newTrack.Id,
			Length: newTrack.Length,
			Name:   newTrack.Name,
		},
	)
	return &newTrack, nil
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

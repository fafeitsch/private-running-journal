package tracks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/shared"
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
	cache    map[string]*Track
	basePath string
}

func New(baseDir string) (*Tracks, error) {
	var trackCache = make(map[string]*Track)
	result := Tracks{cache: trackCache, basePath: filepath.Join(baseDir, "tracks")}
	err := os.MkdirAll(result.basePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create tracks directory: %v", err)
	}

	err = filepath.WalkDir(
		result.basePath, func(path string, info os.DirEntry, err error) error {
			if err != nil {
				log.Printf("skipping directory \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			if info.IsDir() || info.Name() != "info.json" {
				return nil
			}
			parent := strings.Replace(path, "/"+info.Name(), "", 1)
			relativePath := strings.Replace(parent, result.basePath+"/", "", 1)
			log.Printf("relative path: %s", relativePath)
			track, err := result.readTrack(parent, relativePath)
			if err != nil {
				log.Printf("skipping track \"%s\" because an error occurred: %v", path, err)
				return filepath.SkipDir
			}
			result.cache[track.Id] = &track
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not read tracks: %v", err)
	}
	shared.RegisterHandler(
		"journal entry changed", func(data ...any) {
			old := data[0].(shared.JournalEntry)
			nevv := data[1].(shared.JournalEntry)
			if oldTrack, ok := trackCache[old.TrackId]; ok {
				oldTrack.Usages = oldTrack.Usages - 1
			}
			if newTrack, ok := trackCache[nevv.TrackId]; ok {
				newTrack.Usages = newTrack.Usages + 1
			}
		},
	)
	sendInitEvent(trackCache)
	return &result, nil
}

func sendInitEvent(trackCache map[string]*Track) {
	list := make([]shared.Track, 0)
	for _, track := range trackCache {
		list = append(
			list, shared.Track{
				Id:     track.Id,
				Length: track.Length,
				Name:   track.Name,
			},
		)
	}
	shared.Send("tracks initialized", list)
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

	if err != nil {
		return Track{}, err
	}

	hierarchy := strings.Split(relativePath, string(os.PathSeparator))
	return Track{
		Id:        relativePath,
		Length:    length,
		Name:      baseDescriptor.Name,
		Hierarchy: hierarchy[:len(hierarchy)-1],
	}, nil
}

func (t *Tracks) Tracks() []Track {
	result := make([]Track, 0)
	for _, track := range t.cache {
		result = append(result, *track)
	}
	return result
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
	existing, ok := t.cache[track.Id]
	if !ok {
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
	err = os.WriteFile(filepath.Join(trackDirectory, "track.gpx"), writer.Bytes(), 0644)
	if err != nil {
		return nil, err
	}
	existingTrack, err := t.readTrack(trackDirectory, track.Id)
	shared.Send(
		"track upserted", shared.Track{
			Id:     existingTrack.Id,
			Length: existingTrack.Length,
			Name:   existingTrack.Name,
		},
	)
	return &existingTrack, err
}

type CreateTrack struct {
	Name   string `json:"name"`
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
	gpxPayload := gpx.GPX{Trk: []*gpx.TrkType{{}}}
	writer := bytes.Buffer{}
	_ = gpxPayload.WriteIndent(bufio.NewWriter(&writer), "  ", "  ")
	err = os.WriteFile(filepath.Join(trackPath, "track.gpx"), writer.Bytes(), 0644)
	if err != nil {
		return nil, err
	}
	id := strings.Replace(trackPath, t.basePath+"/", "", 1)
	newTrack, err := t.readTrack(trackPath, id)
	if err != nil {
		return nil, err
	}
	t.cache[newTrack.Id] = &newTrack
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
	track, ok := t.cache[id]
	if !ok {
		return nil
	}
	err := os.RemoveAll(filepath.Join(t.basePath, track.Id))
	fmt.Printf("%s map: %v", id, t.cache)
	fmt.Printf("map: %v", t.cache)
	shared.Send("track deleted", id)
	t.deleteEmptyDirectories(id)
	return err
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

type GpxData struct {
	Waypoints       []Coordinates    `json:"waypoints"`
	DistanceMarkers []DistanceMarker `json:"distanceMarkers"`
}

func (t *Tracks) GetGpxData(id string) (GpxData, error) {
	gpxPath := filepath.Join(filepath.Join(t.basePath, id), "track.gpx")
	coordinates, _, err := readGpx(gpxPath)
	distanceMarkers := distanceMarkers(coordinates, 1000)
	return GpxData{Waypoints: coordinates, DistanceMarkers: distanceMarkers}, err
}

type PolylineProps struct {
	Length          int              `json:"length"`
	DistanceMarkers []DistanceMarker `json:"distanceMarkers"`
}

func ComputePolylineProps(coordinates []Coordinates) PolylineProps {
	distanceMarkers := distanceMarkers(coordinates, 1000)
	distance := int(1000 * distance(coordinates))
	return PolylineProps{Length: distance, DistanceMarkers: distanceMarkers}
}

package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/shared"
	"github.com/twpayne/go-gpx"
	"os"
	"path/filepath"
)

type Track struct {
	Id          string   `json:"id"`
	Length      int      `json:"length"`
	Name        string   `json:"name"`
	Variants    []*Track `json:"variants"`
	ParentNames []string `json:"parentNames"`
	Usages      int      `json:"usages"`
}

type Tracks struct {
	cache   map[string]*Track
	baseDir string
}

func New(baseDir string) (*Tracks, error) {
	var trackCache = make(map[string]*Track)
	result := Tracks{cache: trackCache, baseDir: baseDir}
	tracksDir := filepath.Join(baseDir, "tracks")
	baseDirEntries, err := os.ReadDir(tracksDir)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", tracksDir, err)
	}

	for _, baseTrack := range baseDirEntries {
		if !baseTrack.IsDir() {
			continue
		}
		track, err := result.readTrack(filepath.Join(tracksDir, baseTrack.Name()), baseTrack.Name(), []string{})
		if err != nil {
			fmt.Printf("could not read %s, skipping it: %v", tracksDir, err)
			continue
		}
		trackCache[track.Id] = &track
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
				Id:          track.Id,
				Length:      track.Length,
				Name:        track.Name,
				ParentNames: track.ParentNames,
			},
		)
	}
	shared.Send("tracks initialized", list)
}

func (t Tracks) readTrack(path string, relativePath string, parentNames []string) (Track, error) {
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

	subFiles, err := os.ReadDir(path)
	variants := make([]*Track, 0, 0)
	if err != nil {
		return Track{}, err
	}

	for _, subFile := range subFiles {
		if !subFile.IsDir() {
			continue
		}
		parents := append(parentNames, baseDescriptor.Name)
		variant, err := t.readTrack(
			filepath.Join(path, subFile.Name()), filepath.Join(relativePath, subFile.Name()), parents,
		)
		if err != nil {
			return Track{}, fmt.Errorf("could not read variant path %s of %s: %e", subFile.Name(), variant, err)
		}
		t.cache[variant.Id] = &variant
		variants = append(variants, &variant)
	}
	return Track{
		Id:          relativePath,
		Length:      length,
		Name:        baseDescriptor.Name,
		Variants:    variants,
		ParentNames: parentNames,
	}, nil
}

func (t Tracks) RootTracks() []Track {
	result := make([]Track, 0, 0)
	for key, value := range t.cache {
		v := t.cache[key]
		if len(value.ParentNames) == 0 {
			result = append(result, *v)
		}
	}
	return result
}

type SaveTrack struct {
	Id        string        `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Parents   []string      `json:"parents,omitempty"`
	Waypoints []Coordinates `json:"waypoints,omitempty"`
}

func (t Tracks) SaveTrack(track SaveTrack) error {
	return nil
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

func GetGpxData(baseDirectory string, id string) (GpxData, error) {
	path := filepath.Join(baseDirectory, id, "track.gpx")
	coordinates, _, err := readGpx(path)
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

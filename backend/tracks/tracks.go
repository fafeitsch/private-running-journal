package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/twpayne/go-gpx"
	"os"
	"path/filepath"
)

type Track struct {
	Id       string  `json:"id"`
	Length   int     `json:"length"`
	Name     string  `json:"name"`
	Variants []Track `json:"variants"`
}

func GetTracks(baseDir string) ([]Track, error) {
	tracksDir := filepath.Join(baseDir, "tracks")
	baseDirEntries, err := os.ReadDir(tracksDir)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", tracksDir, err)
	}

	result := make([]Track, 0, len(baseDirEntries))
	for _, baseTrack := range baseDirEntries {
		if !baseTrack.IsDir() {
			continue
		}
		track, err := readTrack(filepath.Join(tracksDir, baseTrack.Name()), baseTrack.Name())
		if err != nil {
			fmt.Printf("could not read %s, skipping it: %v", tracksDir, err)
			continue
		}
		result = append(result, track)
	}
	return result, nil
}

func readTrack(path string, relativePath string) (Track, error) {
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
	variants := make([]Track, 0, 0)
	if err != nil {
		return Track{}, err
	}

	for _, subFile := range subFiles {
		if !subFile.IsDir() {
			continue
		}
		variant, err := readTrack(filepath.Join(path, subFile.Name()), filepath.Join(relativePath, subFile.Name()))
		if err != nil {
			return Track{}, fmt.Errorf("could not read variant path %s of %s: %e", subFile.Name(), variant, err)
		}
		variants = append(variants, variant)
	}
	return Track{
		Id:       relativePath,
		Length:   length,
		Name:     baseDescriptor.Name,
		Variants: variants,
	}, nil
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
		length = length + distance(segment.TrkPt)
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

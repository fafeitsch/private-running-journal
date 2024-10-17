package filebased

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/twpayne/go-gpx"
	"os"
	"path/filepath"
	"strings"
)

var tracksDirectory = "tracks"

type trackDescriptor struct {
	Name string `json:"name"`
}

func (s *Service) ReadTrack(path string) (shared.Track, error) {
	descriptorPath := filepath.Join(s.path, "tracks", path, "info.json")
	var baseDescriptor trackDescriptor
	fileContent, err := os.Open(descriptorPath)
	if err != nil {
		return shared.Track{}, err
	}
	err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
	if err != nil {
		return shared.Track{}, err
	}

	gpxPath := filepath.Join(s.path, "tracks", path, "track.gpx")
	waypoints, err := readGpx(gpxPath)
	if err != nil {
		return shared.Track{}, err
	}

	hierarchy := strings.Split(path, string(os.PathSeparator))
	return shared.Track{
		Waypoints: waypoints,
		Id:        path,
		Name:      baseDescriptor.Name,
		Parents:   hierarchy[:len(hierarchy)-1],
	}, nil
}

func readGpx(path string) (shared.Waypoints, error) {
	gpxFileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ocould not read gpx track %s: %v", path, err)
	}
	tracks, err := gpx.Read(bytes.NewReader(gpxFileContent))
	if err != nil {
		return nil, fmt.Errorf("could not parse gpx %s: %v", path, err)
	}
	if len(tracks.Trk) != 1 {
		return nil, fmt.Errorf("%s must contain one track only, but contains %d tracks", path, len(tracks.Trk))
	}
	coordinates := make([]shared.Coordinates, 0, 0)
	for _, segment := range tracks.Trk[0].TrkSeg {
		for _, trkPt := range segment.TrkPt {
			coordinates = append(coordinates, shared.Coordinates{Longitude: trkPt.Lon, Latitude: trkPt.Lat})
		}
	}
	return coordinates, nil
}

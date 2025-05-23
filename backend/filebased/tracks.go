package filebased

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
	"os"
	"path/filepath"
)

var tracksDirectory = "tracks"

type trackDescriptor struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Parents []string `json:"parents"`
	Comment string   `json:"comment"`
}

func (s *Service) ReadAllTracks(consumer func(track shared.Track)) error {
	skipped := make(map[string]error, 0)
	path := filepath.Join(s.path, tracksDirectory)
	return filepath.WalkDir(
		path, func(path string, info os.DirEntry, err error) error {
			if err != nil {
				skipped[path] = err
				return filepath.SkipDir
			}
			if !info.IsDir() || info.Name() == tracksDirectory {
				return nil
			}
			track, err := s.ReadTrack(info.Name())
			if err != nil {
				skipped[path] = err
				return filepath.SkipDir
			}
			consumer(track)
			return nil
		},
	)
}

func (s *Service) ReadTrack(id string) (shared.Track, error) {
	descriptorPath := filepath.Join(s.path, tracksDirectory, id, "info.json")
	var baseDescriptor trackDescriptor
	fileContent, err := os.Open(descriptorPath)
	if err != nil {
		return shared.Track{}, err
	}
	err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
	if err != nil {
		return shared.Track{}, err
	}

	gpxPath := filepath.Join(s.path, "tracks", id, "track.gpx")
	waypoints, err := readGpx(gpxPath)
	if err != nil {
		return shared.Track{}, err
	}

	return shared.Track{
		Waypoints: waypoints,
		Id:        id,
		Name:      baseDescriptor.Name,
		Parents:   baseDescriptor.Parents,
		Comment:   baseDescriptor.Comment,
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

func (s *Service) SaveTrack(track shared.SaveTrack) error {
	trackDirectory := filepath.Join(s.path, "tracks", track.Id)
	err := os.MkdirAll(trackDirectory, 0755)
	if err != nil {
		return fmt.Errorf("could not create track directory %s: %v", trackDirectory, err)
	}
	infoFile := filepath.Join(trackDirectory, "info.json")
	infoPayload, _ := json.Marshal(trackDescriptor{Name: track.Name, Id: track.Id, Parents: track.Parents, Comment: track.Comment})
	err = os.WriteFile(infoFile, infoPayload, 0666)
	if err != nil {
		return fmt.Errorf("could not save base information: %v", err)
	}
	err = writeGpxFile(track.Waypoints, trackDirectory)
	if err != nil {
		return fmt.Errorf("could not save gpx file: %v", err)
	}
	return nil
}

func writeGpxFile(waypoints shared.Waypoints, trackDirectory string) error {
	coords := make([]geom.Coord, 0)
	for _, coordinate := range waypoints {
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

func (s *Service) DeleteTrackDirectory(id string) error {
	err := os.RemoveAll(filepath.Join(s.path, "tracks", id))
	return err
}

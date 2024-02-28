package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/twpayne/go-gpx"
	"log"
	"os"
	"path/filepath"
)

type Track struct {
	Id       string `json:"id"`
	Length   int    `json:"length"`
	BaseName string `json:"baseName"`
	BaseId   string `json:"baseId"`
	Variant  string `json:"variant"`
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
		path := filepath.Join(baseDir, "tracks", baseTrack.Name(), "info.json")
		var baseDescriptor trackDescriptor
		fileContent, err := os.Open(path)
		if err != nil {
			log.Printf("could not open %s, skipping %s", path, baseTrack.Name())
			continue
		}
		err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
		if err != nil {
			log.Printf("could not read baseTrack %s: %v", path, err)
			continue
		}

		variantsPath := filepath.Join(tracksDir, baseTrack.Name())
		variants, err := os.ReadDir(variantsPath)
		if err != nil {
			log.Printf("could not read variants of %s, skipping: %v", baseTrack.Name(), err)
			continue
		}
		for _, variant := range variants {
			if !variant.IsDir() && variant.Name() == "track.gpx" {
				_, length, err := readGpx(filepath.Join(tracksDir, baseTrack.Name(), variant.Name()))
				if err != nil {
					log.Printf("could not read gpx data: %v", err)
					continue
				}
				result = append(
					result, Track{
						BaseName: baseDescriptor.Name,
						Variant:  "",
						BaseId:   baseTrack.Name(),
						Length:   length,
					},
				)
			}
			if !variant.IsDir() {
				continue
			}
			variantTrack, err := readTrackVariant(tracksDir, baseTrack.Name(), baseDescriptor.Name, variant.Name())
			if err != nil {
				log.Printf("could not read track variant %s of base %s: %v", variant.Name(), baseTrack.Name(), err)
				continue
			}
			result = append(result, variantTrack)
		}
	}
	return result, nil
}

type trackDescriptor struct {
	Name string `json:"name"`
}

func readTrackVariant(pathPrefix string, baseId string, baseName string, variantId string) (Track, error) {
	descriptorPath := filepath.Join(pathPrefix, baseId, variantId, "info.json")
	var baseDescriptor trackDescriptor
	fileContent, err := os.Open(descriptorPath)
	if err != nil {
		return Track{}, fmt.Errorf("could not open track descriptor %s: %v", descriptorPath, err)
	}
	err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
	if err != nil {
		return Track{}, fmt.Errorf("could not decode track descriptor %s: %v", descriptorPath, err)
	}
	gpxPath := filepath.Join(pathPrefix, baseId, variantId, "track.gpx")
	_, length, err := readGpx(gpxPath)
	if err != nil {
		return Track{}, fmt.Errorf("could not read gpx data: %v", err)
	}
	return Track{
		Id:       variantId,
		Length:   length,
		BaseId:   baseId,
		BaseName: baseName,
		Variant:  baseDescriptor.Name,
	}, nil
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

func GetGpxData(baseDir string, baseName string, variant string) (GpxData, error) {
	path := filepath.Join(baseDir, "tracks", baseName)
	if variant != "" {
		path = filepath.Join(path, variant)
	}
	path = filepath.Join(path, "track.gpx")
	coordinates, _, err := readGpx(path)
	distanceMarkers := distanceMarkers(coordinates, 1000)
	return GpxData{Waypoints: coordinates, DistanceMarkers: distanceMarkers}, err
}

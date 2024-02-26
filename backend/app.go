package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/twpayne/go-gpx"
	"log"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time now!", name)
}

type JournalListEntry struct {
	Id            string `json:"id"`
	Date          string `json:"date"`
	TrackBaseName string `json:"trackBaseName"`
	TrackVariant  string `json:"trackVariant"`
	Length        int    `json:"length"`
}

type JournalEntry struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Track   Track  `json:"track"`
	Comment string `json:"comment"`
	Time    string `json:"time"`
}

type Track struct {
	Id           string `json:"id"`
	Length       int    `json:"length"`
	BaseName     string `json:"baseName"`
	BaseId       string `json:"baseId"`
	Variant      string `json:"variant"`
	WaypointData string `json:"waypointData"`
}

type TrackListEntryOld struct {
	Name   string  `json:"name"`
	Length float64 `json:"length"`
}

func (a *App) GetJournalListEntries() []JournalListEntry {
	return make([]JournalListEntry, 0)
}

func (a *App) GetJournalEntry() JournalEntry {
	return JournalEntry{}
}

type trackDescriptor struct {
	Name string `json:"name"`
}

func (a *App) GetTracks() []Track {
	baseDirEntries, err := os.ReadDir("appdata/tracks")
	if err != nil {
		log.Printf("could not read appdata/geojson: %v", err)
	}

	result := make([]Track, 0, len(baseDirEntries))
	for _, baseTrack := range baseDirEntries {
		if !baseTrack.IsDir() {
			continue
		}
		path := fmt.Sprintf("./appdata/tracks/%s/info.json", baseTrack.Name())
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

		variantsPath := fmt.Sprintf("./appdata/tracks/%s", baseTrack.Name())
		variants, err := os.ReadDir(variantsPath)
		if err != nil {
			log.Printf("could not read variants of %s, skipping: %v", baseTrack.Name(), err)
			continue
		}
		for _, variant := range variants {
			if !variant.IsDir() && variant.Name() == "track.gpx" {
				data, length, err := a.readGpx(filepath.Join("appdata/tracks", baseTrack.Name(), variant.Name()))
				if err != nil {
					log.Printf("could not read gpx data: %v", err)
					continue
				}
				result = append(
					result,
					Track{
						BaseName:     baseDescriptor.Name,
						Variant:      "",
						BaseId:       baseTrack.Name(),
						Length:       length,
						WaypointData: data,
					},
				)
			}
			if !variant.IsDir() {
				continue
			}
			variantTrack, err := a.readTrackVariant(
				"appdata/tracks", baseTrack.Name(), baseDescriptor.Name, variant.Name(),
			)
			if err != nil {
				log.Printf("could not read track variant %s of base %s: %v", variant.Name(), baseTrack.Name(), err)
				continue
			}
			result = append(result, variantTrack)
		}
	}
	return result
}

func (a *App) readTrackVariant(pathPrefix string, baseId string, baseName string, variantId string) (Track, error) {
	descriptorPath := fmt.Sprintf("%s/%s/%s/info.json", pathPrefix, baseId, variantId)
	var baseDescriptor trackDescriptor
	fileContent, err := os.Open(descriptorPath)
	if err != nil {
		return Track{}, fmt.Errorf("could not open track descriptor %s: %v", descriptorPath, err)
	}
	err = json.NewDecoder(fileContent).Decode(&baseDescriptor)
	if err != nil {
		return Track{}, fmt.Errorf("could not decode track descriptor %s: %v", descriptorPath, err)
	}
	gpxPath := fmt.Sprintf("%s/%s/%s/track.gpx", pathPrefix, baseId, variantId)
	data, length, err := a.readGpx(gpxPath)
	if err != nil {
		return Track{}, fmt.Errorf("could not read gpx data: %v", err)
	}
	return Track{
		Id:           variantId,
		Length:       length,
		BaseId:       baseId,
		BaseName:     baseName,
		Variant:      baseDescriptor.Name,
		WaypointData: data,
	}, nil
}

func (a *App) readGpx(path string) (string, int, error) {
	gpxFileContent, err := os.ReadFile(path)
	if err != nil {
		return "", 0, fmt.Errorf("ocould not read gpx track %s: %v", path, err)
	}
	tracks, err := gpx.Read(bytes.NewReader(gpxFileContent))
	if err != nil {
		return "", 0, fmt.Errorf("could not parse gpx %s: %v", path, err)
	}
	if len(tracks.Trk) != 1 {
		return "", 0, fmt.Errorf("%s must contain one track only, but contains %d tracks", path, len(tracks.Trk))
	}
	length := 0.0
	for _, segment := range tracks.Trk[0].TrkSeg {
		length = length + distance(segment.TrkPt)
	}
	return string(gpxFileContent), int(1000 * length), nil
}

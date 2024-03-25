package settings

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/local-track-journal/backend/shared"
	"os"
	"path/filepath"
)

type MapSettings struct {
	TileServer  string     `json:"tileServer"`
	Attribution string     `json:"attribution"`
	CacheTiles  bool       `json:"cacheTiles"`
	ZoomLevel   int        `json:"zoomLevel"`
	Center      [2]float64 `json:"center"`
}

type AppSettings struct {
	MapSettings MapSettings `json:"mapSettings"`
	HttpPort    int         `json:"httpPort"`
}

type Settings struct {
	settingsFile string
	appSettings  AppSettings
}

func New(baseDirectory string) (*Settings, error) {
	result := Settings{settingsFile: filepath.Join(baseDirectory, "settings.json")}
	err := result.initSettings()
	return &result, err
}

func (s *Settings) initSettings() error {
	s.appSettings = AppSettings{
		MapSettings: MapSettings{
			TileServer:  "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
			Attribution: "&copy; <a href=\"http://www.openstreetmap.org/copyright\">OpenStreetMap</a>",
			CacheTiles:  false,
			ZoomLevel:   6,
			Center:      [2]float64{51.330, 10.453},
		},
		HttpPort: 47836,
	}
	_, err := os.Stat(s.settingsFile)
	if nil == err {
		content, err := os.ReadFile(s.settingsFile)
		if err != nil {
			return err
		}
		return json.Unmarshal(content, &s.appSettings)
	}
	payload, _ := json.MarshalIndent(s.appSettings, "", "  ")
	return os.WriteFile(s.settingsFile, payload, 0664)
}

func (s *Settings) SaveSettings(settings AppSettings) error {
	payload, _ := json.MarshalIndent(settings, "", "  ")
	err := os.WriteFile(s.settingsFile, payload, 0644)
	if err != nil {
		return fmt.Errorf("could not save settings: %v", err)
	}
	if settings.MapSettings.TileServer != s.appSettings.MapSettings.TileServer {
		shared.Send("tile-server-changed", settings.MapSettings.TileServer)
	}
	if settings.MapSettings.CacheTiles != s.appSettings.MapSettings.CacheTiles {
		shared.Send("tile-server-cache-enabled-changed", settings.MapSettings.CacheTiles)
	}
	s.appSettings = settings
	return nil
}

func (s *Settings) AppSettings() AppSettings {
	return s.appSettings
}

func (s *Settings) MapSettings() MapSettings {
	return s.appSettings.MapSettings
}

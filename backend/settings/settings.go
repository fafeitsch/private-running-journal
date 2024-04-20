package settings

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
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
	Language    string      `json:"language"`
	GitSettings GitSettings `json:"gitSettings"`
}

type GitSettings struct {
	Enabled         bool `json:"enabled"`
	PushAfterCommit bool `json:"pushAfterCommit"`
	PullOnStartUp   bool `json:"pullOnStartUp"`
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
		Language: "en",
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
		shared.Send("tile-server-cache-Enabled-changed", settings.MapSettings.CacheTiles)
	}
	if settings.GitSettings.Enabled != s.appSettings.GitSettings.Enabled {
		shared.Send("git enablement changed", settings.GitSettings.Enabled)
	}
	if settings.GitSettings.PushAfterCommit != s.appSettings.GitSettings.PushAfterCommit {
		shared.Send("git push changed", settings.GitSettings.PushAfterCommit)
	}
	shared.Send("settings changed")
	s.appSettings = settings
	return nil
}

func (s *Settings) AppSettings() AppSettings {
	return s.appSettings
}

func (s *Settings) MapSettings() MapSettings {
	return s.appSettings.MapSettings
}

func (s *Settings) GitSettings() GitSettings { return s.appSettings.GitSettings }

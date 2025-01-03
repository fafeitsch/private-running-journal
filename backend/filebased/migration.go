package filebased

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const currentVersion = 2

type fileVersion struct {
	Version int `json:"version"`
}

var migrations = map[int]migrator{1: insertIdsAndRestructureFiles}

func (s *Service) Migrate() error {
	version, err := s.loadCurrentFileVersion()
	if err != nil {
		return err
	}
	log.Printf("found file version %v, current version is %v", version.Version, currentVersion)
	for v := version.Version; v < currentVersion; v = v + 1 {
		log.Printf("starting migration from version %d to version %d", v, v+1)
		err := migrations[v](s.path)
		if err != nil {
			return fmt.Errorf("error migrating from version %d to version %d: %v", v, v+1, err)
		}
		log.Printf("finished migration from version %d to version %d", v, v+1)
		shared.SendEvent(shared.MigrationEvent{
			OldVersion: v,
			NewVersion: v + 1,
		})
		version.Version = v + 1
		payload, _ := json.Marshal(version)
		err = os.WriteFile(filepath.Join(s.path, "fileVersion.json"), payload, 0644)
		if err != nil {
			return fmt.Errorf("could not write new version to file")
		}
	}
	log.Printf("migration finished successfully")
	return nil
}

func (s *Service) loadCurrentFileVersion() (*fileVersion, error) {
	_, err := os.Stat(filepath.Join(s.path, "fileVersion.json"))
	var version fileVersion
	if os.IsNotExist(err) {
		if currentVersion == 2 {
			version = fileVersion{Version: 1}
		} else {
			version = fileVersion{Version: currentVersion}
		}
	} else {
		file, err := os.Open(filepath.Join(s.path, "fileVersion.json"))
		if err != nil {
			return nil, fmt.Errorf("could not open versions file, but it seems to exist: %v", err)
		}
		defer file.Close()
		err = json.NewDecoder(file).Decode(&version)
		if err != nil {
			return nil, fmt.Errorf("could not parse versions file, but it seems to exist: %v", err)
		}
	}
	return &version, nil
}

type migrator func(directory string) error

func insertIdsAndRestructureFiles(directory string) error {
	tracks := make(map[string]string)
	walkPath := filepath.Join(directory, "tracks")
	_, err := os.Stat(walkPath)
	if os.IsNotExist(err) {
		return nil
	}
	err = filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || info.Name() != "info.json" {
			return nil
		}
		var trackInfo map[string]any
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		id := shared.UniqueId()
		_ = json.NewDecoder(file).Decode(&trackInfo)
		trackInfo["id"] = id
		oldId := strings.TrimPrefix(filepath.Dir(path), walkPath)
		tracks[oldId[1:]] = id
		parents := strings.Split(oldId, string(filepath.Separator))
		trackInfo["parents"] = parents[1 : len(parents)-1]
		payload, _ := json.Marshal(trackInfo)
		err = os.WriteFile(path, payload, 0644)
		if err != nil {
			return err
		}
		err = os.Rename(filepath.Dir(path), filepath.Join(directory, "tracks", id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	counter, err := removeEmptyDirs(walkPath)
	for counter > 0 && err == nil {
		counter, err = removeEmptyDirs(walkPath)
	}
	if err != nil {
		return err
	}
	walkPath = filepath.Join(directory, "journal")
	err = filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || info.Name() != "entry.json" {
			return nil
		}
		var entryInfo map[string]any
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		id := shared.UniqueId()
		_ = json.NewDecoder(file).Decode(&entryInfo)
		entryInfo["id"] = id
		oldId := strings.TrimPrefix(filepath.Dir(path), walkPath)
		tracks[id] = oldId
		parents := strings.Split(oldId, string(filepath.Separator))
		entryInfo["date"] = fmt.Sprintf("%s-%s-%s", parents[1], parents[2], parents[3][0:2])
		entryInfo["track"] = tracks[entryInfo["track"].(string)]
		payload, _ := json.Marshal(entryInfo)
		err = os.WriteFile(path, payload, 0644)
		if err != nil {
			return err
		}
		err = os.MkdirAll(filepath.Join(directory, "journal", id[0:2]), 0755)
		if err != nil {
			return err
		}
		err = os.Rename(filepath.Dir(path), filepath.Join(directory, "journal", id[0:2], id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	counter, err = removeEmptyDirs(walkPath)
	for counter > 0 && err == nil {
		counter, err = removeEmptyDirs(walkPath)
	}
	return err
}

func removeEmptyDirs(root string) (int, error) {
	counter := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || !info.IsDir() {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}
		if len(entries) > 0 {
			return nil
		}
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to remove directory %s: %v", path, err)
		}
		counter = counter + 1

		return nil
	})
	return counter, err
}

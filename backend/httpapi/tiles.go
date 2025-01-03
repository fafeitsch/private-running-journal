package httpapi

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tilesDirMutex sync.Mutex

type TileServer struct {
	baseDir      string
	url          string
	cacheEnabled bool
}

func NewTileServer(baseDir string, url string, cacheEnabled bool) *TileServer {
	result := TileServer{url: url, baseDir: filepath.Join(baseDir, "tiles"), cacheEnabled: cacheEnabled}
	shared.Listen(shared.TileServerChangedEvent{}, func(k shared.TileServerChangedEvent) {
		result.url = k.NewValue
		_ = os.RemoveAll(result.baseDir)
	})
	shared.Listen(shared.TileServerCacheEnabledEvent{}, func(k shared.TileServerCacheEnabledEvent) {
		result.cacheEnabled = k.NewValue
		if !result.cacheEnabled {
			_ = os.RemoveAll(result.baseDir)
		}
	})
	return &result
}

func (t *TileServer) ServeHTTP(resp http.ResponseWriter, originalRequest *http.Request) {
	parts := strings.Split(originalRequest.URL.String(), "/")
	z := parts[2]
	x := parts[3]
	y := parts[4]
	resp.Header().Set("Content-Type", "image/png")
	if t.readCacheFile(z, x, y, resp) {
		return
	}
	url := strings.Replace(t.url, "{z}", z, 1)
	url = strings.Replace(url, "{x}", x, 1)
	url = strings.Replace(url, "{y}", y, 1)
	client := http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "github.com/fafeitsch/private-running-journal")
	response, err := client.Do(req)
	if err != nil {
		http.Error(resp, "could not get tile from tile server", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(response.Body)
	if t.cacheEnabled {
		go func() {
			tilesDirMutex.Lock()
			defer tilesDirMutex.Unlock()
			_ = os.MkdirAll(filepath.Join(t.baseDir, z, x), os.ModePerm)
			_ = os.WriteFile(filepath.Join(t.baseDir, z, x, y)+".png", body, 0644)
		}()
	}
	resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
	resp.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = resp.Write(body)
}

func (t *TileServer) readCacheFile(z string, x string, y string, resp http.ResponseWriter) bool {
	tilesDirMutex.Lock()
	defer tilesDirMutex.Unlock()
	cachedFile, err := os.Stat(filepath.Join(t.baseDir, z, x, y) + ".png")
	if err != nil || time.Now().Sub(cachedFile.ModTime()) >= 24*time.Hour*180 {
		return false
	}
	tile, err := os.ReadFile(filepath.Join(t.baseDir, z, x, y) + ".png")
	resp.Header().Set("Content-Length", strconv.Itoa(len(tile)))
	resp.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	resp.Header().Set("Pragma", "no-cache")
	resp.Header().Set("Expires", "0")
	if err == nil {
		_, _ = resp.Write(tile)
	} else {
		fmt.Printf("bigg error %v", err)
	}
	return true
}

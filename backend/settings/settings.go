package settings

type MapSettings struct {
	TileServer  string `json:"tileServer"`
	Attribution string `json:"attribution"`
	CacheTiles  bool   `json:"cacheTiles"`
}

package shared

type Track struct {
	Id          string   `json:"id"`
	Length      int      `json:"length"`
	Name        string   `json:"name"`
	ParentNames []string `json:"parentNames"`
}

type JournalEntry struct {
	TrackId string `json:"trackId"`
}

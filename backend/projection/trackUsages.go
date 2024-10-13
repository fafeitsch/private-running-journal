package projection

import (
	"encoding/json"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"os"
	"path/filepath"
)

func (p *Projection) buildUsages(entries []shared.JournalEntry) error {
	usages := make(map[string]int)
	for _, entry := range entries {
		usages[entry.TrackId] = usages[entry.TrackId] + 1
	}
	payload, _ := json.Marshal(usages)
	return os.WriteFile(filepath.Join(p.directory, "trackUsages.json"), payload, 0644)
}

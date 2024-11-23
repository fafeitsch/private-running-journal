package journalEditor

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
)

type JournalEditor struct {
	fileService *filebased.Service
	trackLookup *projection.TrackIdMap
}

func New(service *filebased.Service, trackLookup *projection.TrackIdMap) *JournalEditor {
	return &JournalEditor{fileService: service, trackLookup: trackLookup}
}
